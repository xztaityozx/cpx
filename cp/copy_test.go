package cp

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vbauerster/mpb"
)

func Test_New(t *testing.T) {
	as := assert.New(t)

	for _, v := range []struct {
		src string
		dst string
	}{
		{src: "/path/to/src", dst: "/path/to/dst"},
	} {
		expect := New(v.src, v.dst)
		as.Equal(expect.src, v.src)
		as.Equal(expect.dst, v.dst)
	}
}

func Test_WithProgressBar(t *testing.T) {
	as := assert.New(t)
	parent := mpb.New()

	t.Run("failed open src", func(t *testing.T) {
		cp := New("__", "dst")
		as.Error(cp.WithProgressBar(parent))
	})

	t.Run("ok", func(t *testing.T) {
		tmp := os.TempDir()
		dir := filepath.Join(tmp, "cpx")
		os.MkdirAll(dir, 0755)
		path := filepath.Join(dir, "file")

		w, _ := os.Create(path)
		io.WriteString(w, "test")

		cp := New(path, "")
		as.NoError(cp.WithProgressBar(parent))

		os.RemoveAll(dir)
	})

}
