package ftp

import (
	"context"
	"fmt"
	"io"

	"github.com/jlaffaye/ftp"
)

type Client interface {
	Upload(ctx context.Context, addr, user, pass string, r io.Reader, remoteDir, remoteFilename string) error
}

type ftpClient struct{}

func (s *ftpClient) Upload(ctx context.Context, addr, user, pass string, r io.Reader, remoteDir, remoteFilename string) error {
	conn, err := ftp.Dial(addr, ftp.DialWithContext(ctx))
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer func() { _ = conn.Quit() }()

	if err := conn.Login(user, pass); err != nil {
		return fmt.Errorf("login: %w", err)
	}

	if err := conn.ChangeDir(remoteDir); err != nil {
		// TODO(sgc): Call MakeDir and try again
		return fmt.Errorf("change dir %q: %w", remoteDir, err)
	}

	if err := conn.Stor(remoteFilename, r); err != nil {
		return fmt.Errorf("store: %w", err)
	}

	return nil
}
