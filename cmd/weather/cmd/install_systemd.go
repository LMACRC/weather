//go:build linux && darwin
// +build linux darwin

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/lmacrc/weather"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newInstallSystemd())
}

func newInstallSystemd() *cobra.Command {
	var flags = struct {
		force bool
	}{}

	cmd := &cobra.Command{
		Use:   "install-systemd",
		Short: "Install weather as a systemd service",
		RunE: func(cmd *cobra.Command, args []string) error {
			var is installSystemd

			if err := is.hasSystemd(); err != nil {
				return err
			}

			if !flags.force && is.checkFileExists("/lib/systemd/system/weather.service") {
				fmt.Println("Service installed")
				return nil
			}

			return is.copyService()
		},
	}

	cmd.Flags().BoolVarP(&flags.force, "force", "f", false, "Force overwrite over weather.service file")

	return cmd
}

type installSystemd struct{}

func (i installSystemd) hasSystemd() error {
	if fi, err := os.Stat("/lib/systemd/system"); err != nil {
		return fmt.Errorf("systemd not present")
	} else if !fi.IsDir() {
		return fmt.Errorf("/lib/systemd/system is not a directory")
	}

	return nil
}

// checkFileExists determines if specified file exists
func (i installSystemd) checkFileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

// checkDirExists determines if specified directory exists
func (i installSystemd) checkDirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func (i installSystemd) linkPath(from, to string) error {
	return os.Symlink(from, to)
}

func (i installSystemd) copyService() error {
	// copy service file
	fs, err := weather.Content.Open("etc/weather.service")
	if err != nil {
		return err
	}

	dst, err := os.Create("/lib/systemd/system/weather.service")
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("insufficient permissions to install service")
		}
		return err
	}

	_, err = io.Copy(dst, fs)
	return err
}
