package brotli

import (
	"errors"
	"fmt"
	"os/exec"
)

var (
	ErrNotAvailable = errors.New("brotli command is not in search path")
	ErrCorrupt      = errors.New("corrupt input")
)

// IsAvailable returns true if brotli compression is available.
func IsAvailable() bool {
	_, err := exec.LookPath("brotli")
	return err == nil
}

func TestArchive(path string) error {
	if !IsAvailable() {
		return ErrNotAvailable
	}

	cmd := exec.Command("brotli", "-t", path)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	if cmd.ProcessState.Success() {
		return nil
	}

	return ErrCorrupt
}
