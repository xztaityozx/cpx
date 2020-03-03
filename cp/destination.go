package cp

import (
	"io"
	"os"

	"golang.org/x/xerrors"
)

type Destination struct {
	path string
	w    io.WriteCloser
}

func (d *Destination) Write(p []byte) (n int, err error) {
	return d.w.Write(p)
}

func (d *Destination) Close() error {
	return d.w.Close()
}

// Create Destination struct
func Dst(path string, prompt func(string) bool, force bool) (*Destination, error) {
	var rt Destination
	rt.path = path

	fi, err := os.Stat(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to open destination file(%s): error: %w", path, err)
	}

	return &rt, nil
}
