package brotli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Writer struct {
	pr, pw *os.File
	cmd    *exec.Cmd
}

func NewWriter(w io.Writer) (*Writer, error) {
	if !IsAvailable() {
		return nil, ErrNotAvailable
	}

	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("pipe: %w", err)
	}

	wr := &Writer{
		pr:  pr,
		pw:  pw,
		cmd: exec.Command("brotli", "-Z"),
	}

	wr.cmd.Stdin = pr
	wr.cmd.Stdout = w
	wr.cmd.Stderr = os.Stderr

	if err := wr.cmd.Start(); err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	return wr, nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.pw.Write(p)
}

func (w *Writer) Close() error {
	if w.pw != nil {
		defer func() {
			*w = Writer{}
		}()

		err := w.pw.Close()
		err2 := w.cmd.Wait()

		if err != nil {
			return fmt.Errorf("pipe writer close: %w", err)
		}

		if err2 != nil {
			return fmt.Errorf("brotli: %w", err2)
		}
		_ = w.pr.Close()
	}

	return nil
}
