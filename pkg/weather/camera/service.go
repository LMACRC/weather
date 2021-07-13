package camera

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lmacrc/weather/pkg/filepath/template"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

type Ftp interface {
	Upload(ctx context.Context, dir, filename string, r io.Reader) error
}

type Service struct {
	log       *zap.Logger
	ftp       Ftp
	camera    Capturer
	params    CaptureParams
	localDir  string
	remoteDir string
	filename  *template.Template
	schedule  cron.Schedule
}

func New(log *zap.Logger, vp *viper.Viper, ftp Ftp) (*Service, error) {
	var cfg Config
	if err := vp.UnmarshalKey("camera", &cfg); err != nil {
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
		log:       log.With(zap.String("service", "camera")),
		ftp:       ftp,
		camera:    camera,
		params:    cfg.CaptureParams,
		localDir:  cfg.LocalDir,
		remoteDir: cfg.RemoteDir,
		filename:  template.Must(template.New("file").Parse(cfg.Filename)),
		schedule:  schedule,
	}

	return s, nil
}

func (s Service) Run(ctx context.Context) {
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

	fullPath := filepath.Join(s.localDir, buf.String())

	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	src, err := s.decodeImage(path)
	dst := image.NewRGBA(src.Bounds())
	draw.Copy(dst, image.Pt(0, 0), src, src.Bounds(), draw.Over, nil)

	s.addTimestamp(dst, ts)

	err = jpeg.Encode(file, dst, &jpeg.Options{Quality: 100})
	if err != nil {
		return "", fmt.Errorf("jpeg encode: %w", err)
	}

	return fullPath, nil
}

func (s Service) addTimestamp(img *image.RGBA, ts time.Time) {
	label := ts.Format(time.RFC850)

	bottom := img.Bounds().Max.Y - 10
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{R: 200, G: 100, A: 255}),
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

	s.log.Info("Starting image upload.", zap.String("path", fullPath))
	err = s.ftp.Upload(ctx, s.remoteDir, filepath.Base(fullPath), file)
	if err != nil {
		s.log.Error("Image upload failed.", zap.Error(err))
	} else {
		s.log.Info("Image upload succeeded.")
	}
}
