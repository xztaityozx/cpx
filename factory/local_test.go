package factory

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateLocalCopyers(t *testing.T) {
	as := assert.New(t)
	tmp := filepath.Join(os.TempDir(), "cpx")
	os.MkdirAll(tmp, 0755)

	os.Chdir(tmp)
	files := []string{}
	for i := 0; i < 10; i++ {
		p := filepath.Join(tmp, fmt.Sprint(i))
		ioutil.WriteFile(p, []byte(fmt.Sprint("this is ", i)), 0644)
		files = append(files, p)
	}

	dir1 := filepath.Join(tmp, "d1")
	dir2 := filepath.Join(dir1, "d2")

	os.MkdirAll(dir1, 0755)
	os.MkdirAll(dir2, 0755)

	for i := 0; i < 5; i++ {
		d := dir1
		if i > 2 {
			d = dir2
		}
		p := filepath.Join(d, fmt.Sprint(i))
		ioutil.WriteFile(p, []byte(fmt.Sprint("this is ", i, "(", d, ")")), 0644)
		files = append(files, p)
	}

	t.Run("file only(dst is not exists directory)", func(t *testing.T) {
		dst := filepath.Join(tmp, "dst1")
		res, err := GenerateLocalCopyers(files[:9], []string{dst}, false)
		as.NoError(err)
		as.Equal(len(files[:9]), len(res))
		for i, v := range res {
			as.Equal(files[i], v.Src())
			as.Equal(filepath.Join(dst, fmt.Sprint(i)), v.Dst())
		}

		os.RemoveAll(dst)
	})

	t.Run("file only(dst is exists)", func(t *testing.T) {
		dst := filepath.Join(tmp, "dst2")
		os.MkdirAll(dst, 0755)
		res, err := GenerateLocalCopyers(files[:9], []string{dst}, false)
		as.NoError(err)
		as.Equal(len(files[:9]), len(res))
		for i, v := range res {
			as.Equal(files[i], v.Src())
			as.Equal(filepath.Join(dst, fmt.Sprint(i)), v.Dst())
		}

		os.RemoveAll(dst)
	})

	t.Run("recursive", func(t *testing.T) {
		dst := filepath.Join(tmp, "dst3")
		res, err := GenerateLocalCopyers([]string{tmp}, []string{dst}, true)
		as.NoError(err)
		as.Equal(len(files), len(res))

		for i, v := range res {
			as.Equal(files[i], v.Src())
		}
	})
	os.RemoveAll(tmp)
}
