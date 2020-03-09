package factory

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/xztaityozx/cpx/cp"
	"golang.org/x/xerrors"
)

type Factory struct {
	prompt func(string) bool
	force  bool
}

type Copies []cp.Copy

func New(prompt func(string) bool, force bool) Factory { return Factory{prompt: prompt, force: force} }

func (f Factory) openLocal(src, dst string) (*cp.Source, *cp.Destination, error) {
	s, err := cp.File(src)
	if err != nil {
		return nil, nil, err
	}
	d, err := cp.Dst(dst, f.prompt, f.force)

	return s, d, err
}

// NOT tested
func (f Factory) HttpGet(url, dst string) (cp.Copy, error) {
	src, err := cp.HttpGet(url)
	if err != nil {
		return cp.Copy{}, err
	}
	r, _ := http.NewRequest("GET", url, nil)
	filename := path.Base(r.URL.Path)
	dst = filepath.Join(dst, filename)
	d, err := cp.Dst(dst, f.prompt, f.force)
	if err != nil {
		return cp.Copy{}, err
	}
	return cp.New(src, d), nil
}

func (f Factory) Directory(src, dst string, recursive bool) (rt *Copies, err error) {
	src, dst = filepath.Clean(src), filepath.Clean(dst)
	srcFi, err := os.Stat(src)
	if err != nil {
		return nil, err
	}
	if !srcFi.IsDir() {
		return nil, xerrors.New("source is not directory")
	}

	dstFi, err := os.Stat(src)
	if err != nil {
		if err := os.MkdirAll(dst, srcFi.Mode().Perm()); err != nil {
			return nil, err
		}
		dstFi, _ = os.Stat(dst)
	}

	if !dstFi.IsDir() {
		return nil, xerrors.New("destination is not directory")
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		s, d := filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())

	}

}
