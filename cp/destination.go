package cp

import (
	"fmt"
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
	if err == nil {
		if fi.IsDir() {
			return nil, xerrors.Errorf("%s is directory", path)
		} else if !fi.Mode().IsRegular() {
			return nil, xerrors.Errorf("%s is not regular file", path)
		}

		if !force && !prompt(fmt.Sprintf("%s is already exists. overwrite it?", path)) {
			return nil, xerrors.Errorf("canceled by user")
		}
	}

	rt.w, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}
