package cp

import (
	"io"
	"net/http"
	"os"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"golang.org/x/xerrors"
)

type Source struct {
	r    io.ReadCloser
	pb   func() *mpb.Bar
	size int64
	path string
}

func (s *Source) Read(p []byte) (n int, err error) {
	return s.r.Read(p)
}

func (s *Source) Close() error {
	return s.r.Close()
}

// Create Source struct from file
func File(path string) (*Source, error) {
	fp, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, xerrors.Errorf("failed to open source file: %s\nerror-> %w", path, err)
	}
	fi, _ := fp.Stat()
	return &Source{r: fp, size: fi.Size(), path: path}, nil
}

// Create Source struct from http.get
func HttpGet(url string) (*Source, error) {
	res, err := http.Get(url)
	return &Source{r: res.Body, size: res.ContentLength, path: url}, err
}

// enable
func (s *Source) WithProgressBar(parent *mpb.Progress) {
	s.pb = func() *mpb.Bar {
		p := parent.AddBar(s.size,
			mpb.BarStyle("[=>-|"),
			mpb.PrependDecorators(
				decor.CountersKibiByte("% .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.EwmaETA(decor.ET_STYLE_GO, 90),
				decor.Name(" ] "),
				decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
			))
		s.r = p.ProxyReader(s.r)
		return p
	}
}
