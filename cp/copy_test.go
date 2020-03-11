package cp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Glob(t *testing.T) {
	as := assert.New(t)
	tmp := filepath.Join(os.TempDir(), "cpx")
	as.NoError(os.MkdirAll(tmp, 0755), "failed to create test directory")

	defer os.RemoveAll(tmp)

	for i := 0; i < 10; i++ {
		p := filepath.Join(tmp, fmt.Sprintf("f%d", i))
		as.NoError(ioutil.WriteFile(p, []byte(fmt.Sprintf("this is test file%d", i)), 0644))
	}

	d1 := filepath.Join(tmp, "d1")
	d2 := filepath.Join(d1, "d2")
	as.NoError(os.MkdirAll(d2, 0755), "failed to create test directory")

	for i := 0; i < 3; i++ {
		x, y := filepath.Join(d1, fmt.Sprintf("f%d", i)), filepath.Join(d2, fmt.Sprintf("f%d", i))

		as.NoError(ioutil.WriteFile(x, []byte(fmt.Sprintf("this is test file%d(d1)", i)), 0644))
		as.NoError(ioutil.WriteFile(y, []byte(fmt.Sprintf("this is test file%d(d2)", i)), 0644))
	}

	type pair struct {
		s string
		d string
	}

	dst1 := filepath.Join(tmp, "dst1")
	//dst2 := filepath.Join(tmp, "dst2")

	yes := func(s string) bool { return true }
	//no := func(s string) bool { return false }

	data := []struct {
		glob   []string
		dst    []string
		expect []pair
		skip   []string
		r      bool
		throw  bool
		prompt func(s string) bool
	}{
		{glob: []string{d1}, dst: []string{dst1}, r: false, throw: false, expect: []pair{
			{s: filepath.Join(d1, "f0"), d: filepath.Join(dst1, "f0")},
			{s: filepath.Join(d1, "f1"), d: filepath.Join(dst1, "f1")},
			{s: filepath.Join(d1, "f2"), d: filepath.Join(dst1, "f2")},
		}, skip: []string{d2}, prompt: yes},
	}

	for _, v := range data {
		res, sk, err := Glob(v.glob, v.dst, v.r, v.prompt)
		if v.throw {
			as.Nil(res)
			as.Nil(sk)
			as.Error(err)
		} else {
			var pairs []pair
			for _, c := range res {
				pairs = append(pairs, pair{s: c.src.path, d: c.dst.path})
			}

			as.Equal(v.expect, pairs)
			as.Equal(v.skip, sk)
		}
	}

}
