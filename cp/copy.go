package cp

import (
	"io"
	"os"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"golang.org/x/xerrors"
)

type Copy struct {
	src string
	dst string
	pb  func() *mpb.Bar
}

// create Copy struct
func New(src, dst string) Copy {
	return Copy{src: src, dst: dst}
}

func (cp Copy) open() (src, dst *os.File, err error) {
	src, err = os.OpenFile(cp.src, os.O_RDONLY, 0644)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to open src file:(%s)\ninner error: %v", cp.src, err)
	}

	if cp.ExistsDst() {
		if err := os.Remove(cp.dst); err != nil {
			return nil, nil, xerrors.Errorf("failed to remove previous dst file(%s): error:%v", cp.dst, err)
		}
	}
	dst, err = os.OpenFile(cp.dst, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to open dst file:(%s)\ninner error: %v", cp.dst, err)
	}

	return
}

// enable progress bar
func (cp *Copy) WithProgressBar(parent *mpb.Progress) error {
	src, err := os.Stat(cp.src)
	if err != nil {
		return xerrors.Errorf("failed to open src file:(%s)", cp.src)
	}

	cp.pb = func() *mpb.Bar {
		return parent.AddBar(src.Size(),
			mpb.BarStyle("[=>-|"),
			mpb.PrependDecorators(
				decor.CountersKibiByte("% .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.EwmaETA(decor.ET_STYLE_GO, 90),
				decor.Name(" ] "),
				decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
			))
	}
	return nil
}

// execute copy
func (cp Copy) Copy() error {
	src, dst, err := cp.open()
	if err != nil {
		return err
	}

	if st, _ := src.Stat(); !st.Mode().IsRegular() {
		return xerrors.Errorf("%s is not regular file", cp.src)
	}
	if st, _ := dst.Stat(); !st.Mode().IsRegular() {
		return xerrors.Errorf("%s is not regular file", cp.dst)
	}

	r := func() io.ReadCloser {
		if cp.pb != nil {
			return cp.pb().ProxyReader(src)
		}
		return src
	}()
	defer func() { r.Close(); dst.Close() }()

	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	err = dst.Sync()
	if err != nil {
		return err
	}
	fi, err := os.Stat(cp.src)
	if err != nil {
		return err
	}
	return os.Chmod(cp.dst, fi.Mode())
}

func (cp Copy) ExistsDst() bool {
	_, err := os.Stat(cp.dst)
	return err == nil
}

func (cp Copy) Src() string {
	return cp.src
}
func (cp Copy) Dst() string {
	return cp.dst
}
