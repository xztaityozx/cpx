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

	data := []struct {
		glob   string
		expect []string
		r      bool
		throw  bool
	}{
		{glob: d1, r: false, throw: false, expect: []string{
			filepath.Join(d1, "f0"),
			filepath.Join(d1, "f1"),
			filepath.Join(d1, "f2"),
			d1,
		}},
		{glob: d1, r: true, throw: false, expect: []string{
			filepath.Join(d1, "f0"),
			filepath.Join(d1, "f1"),
			filepath.Join(d1, "f2"),
			filepath.Join(d2, "f0"),
			filepath.Join(d2, "f1"),
			filepath.Join(d2, "f2"),
			d2, d1,
		}},
		{glob: d2, r: true, throw: false, expect: []string{
			filepath.Join(d2, "f0"),
			filepath.Join(d2, "f1"),
			filepath.Join(d2, "f2"),
			d2,
		}},
		{glob: filepath.Join(d1, "*"), r: false, throw: false, expect: []string{
			filepath.Join(d1, "f0"),
			filepath.Join(d1, "f1"),
			filepath.Join(d1, "f2"),
		}},
	}

	for _, v := range data {
		res, err := Glob(v.glob, v.r)
		if v.throw {
			as.Error(err, v.glob)
			as.Nil(res, v.glob)
		} else {
			as.NoError(err, v.glob)
			actual := []string{}
			for _, fe := range res {
				actual = append(actual, filepath.Join(fe.base, fe.path))
			}
			as.ElementsMatch(v.expect, actual, "glob: %s, r: %v", v.glob, v.r)
		}
	}

}
