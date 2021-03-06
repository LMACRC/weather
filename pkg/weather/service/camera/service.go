package camera

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lmacrc/weather/pkg/filepath/template"
	"github.com/lmacrc/weather/pkg/weather/service"
	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

type Service struct {
	log          *zap.Logger
	ftp          service.Ftp
	camera       Capturer
	params       CaptureParams
	outputParams OutputParams
	localDir     string
	remoteDir    string
	filename     *template.Template
	schedule     cron.Schedule
}

func New(log *zap.Logger, vp *viper.Viper, ftp service.Ftp) (*Service, error) {
	var cfg Config
	if err := vp.UnmarshalKey("camera", &cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc())); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	schedule, err := cron.ParseStandard(cfg.Cron)
	if err != nil {
		return nil, fmt.Errorf("parsing cron: %w", err)
	}

	var camera Capturer
	if driverFn := drivers[cfg.Driver]; driverFn == nil {
		return nil, fmt.Errorf("invalid camera driver %q: expect [%s]", cfg.Driver, strings.Join(driverList(), ", "))
	} else if camera, err = driverFn(vp); err != nil {
		return nil, fmt.Errorf("camera driver create: %w", err)
	}

	s := &Service{
		log:          log.With(zap.String("service", "camera")),
		ftp:          ftp,
		camera:       camera,
		params:       cfg.CaptureParams,
		outputParams: cfg.OutputParams,
		localDir:     cfg.LocalDir,
		remoteDir:    cfg.RemoteDir,
		filename:     template.Must(template.New("file").Parse(cfg.Filename)),
		schedule:     schedule,
	}

	return s, nil
}

func (s Service) Run(ctx context.Context) {
	s.log.Info("Starting.")

	for {
		ts := time.Now()
		next := s.schedule.Next(ts)
		sleep := next.Sub(ts)
		s.log.Info("Next upload scheduled.", zap.Time("time", next), zap.Duration("wait_time", sleep))

		select {
		case <-ctx.Done():
			s.log.Info("Shutting down.")
			return

		case <-time.After(sleep):
			s.processNextImage(ctx, next)
		}
	}
}

func (s Service) decodeImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image decode: %w", err)
	}
	return img, err
}

func (s Service) CaptureImage(_ context.Context, ts time.Time) (string, error) {
	path, err := s.camera.Capture(s.params)
	if err != nil {
		return "", fmt.Errorf("capture: %w", err)
	}
	defer func() { _ = os.Remove(path) }()

	var buf bytes.Buffer
	err = s.filename.Execute(&buf, map[string]interface{}{
		"Now": ts,
	})
	if err != nil {
		return "", fmt.Errorf("filename template: %w", err)
	}

	dstPath := filepath.Join(s.localDir, buf.String())
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer func() { _ = dstFile.Close() }()

	src, err := s.decodeImage(path)
	if err != nil {
		return "", fmt.Errorf("decode image: %w", err)
	}

	dst := image.NewRGBA(src.Bounds())
	draw.Copy(dst, image.Pt(0, 0), src, src.Bounds(), draw.Over, nil)

	s.addTimestamp(dst, ts)

	err = jpeg.Encode(dstFile, dst, &jpeg.Options{Quality: 100})
	if err != nil {
		return "", fmt.Errorf("jpeg encode: %w", err)
	}

	return dstPath, nil
}

func (s Service) addTimestamp(img *image.RGBA, ts time.Time) {
	label := ts.Format(time.RFC850)

	bottom := img.Bounds().Max.Y - 10
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(s.outputParams.TextColor.ToRGBA()),
		Face: inconsolata.Bold8x16,
		Dot:  fixed.P(10, bottom),
	}

	d.DrawString(label)
}

func (s Service) processNextImage(ctx context.Context, ts time.Time) {
	s.log.Info("Generating webcam image.")

	fullPath, err := s.CaptureImage(ctx, ts)
	if err != nil {
		s.log.Error("Failed to capture image.", zap.Error(err))
		return
	}

	file, err := os.Open(fullPath)
	if err != nil {
		s.log.Error("Failed to open captured image file.", zap.String("path", fullPath), zap.Error(err))
		return
	}
	defer func() { _ = file.Close() }()

	if s.ftp == nil {
		// FTP service is disabled
		return
	}

	s.log.Info("Enqueue image upload.", zap.String("path", fullPath))
	err = s.ftp.Enqueue(service.FtpRequest{
		LocalPath:      fullPath,
		RemoteDir:      s.remoteDir,
		RemoteFilename: filepath.Base(fullPath),
		RemoveLocal:    false,
	})
	if err != nil {
		s.log.Error("Failed to queue image upload.", zap.Error(err))
	}
}
