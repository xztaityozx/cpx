package ff

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_YesNo(t *testing.T) {
	as := assert.New(t)
	var ff FuzzyFinder = FuzzyFinder{
		Command: "head",
	}
	ff.Options = []string{"-n1"}

	t.Run("yes", func(t *testing.T) {
		as.True(ff.YesNo(""))
	})

	t.Run("no", func(t *testing.T) {
		ff.Command = "/usr/bin/tail"
		as.False(ff.YesNo(""))
	})

	t.Run("no-2", func(t *testing.T) {
		ff.Command = "/bin/cat"
		ff.Options = []string{}
		as.False(ff.YesNo(""))
	})

	t.Run("no-3", func(t *testing.T) {
		ff.Command = "echo"
		as.False(ff.YesNo(""))
	})
}

func Test_GetPathes(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "cpx")
	os.MkdirAll(tmp, 0755)
	os.Chdir(tmp)

	{
		for i := 0; i < 10; i++ {
			p := filepath.Join(tmp, fmt.Sprint(i))
			ioutil.WriteFile(p, []byte("this is test file"), 0644)
		}

		for i := 1; i <= 3; i++ {
			os.MkdirAll(filepath.Join(tmp, fmt.Sprintf("d%d", i)), 0755)
			for j := 0; j < 3; j++ {
				ioutil.WriteFile(filepath.Join(tmp, fmt.Sprintf("d%d", i), fmt.Sprintf("f%d", j)), []byte("this is f"), 0644)
			}
		}
	}

	defer os.RemoveAll(tmp)

	as := assert.New(t)
	ff := FuzzyFinder{Command: "head"}

	data := []struct {
		opt     string
		glob    string
		expects []string
		throw   bool
	}{
		{opt: "-n13", glob: filepath.Join(tmp, "./??"), expects: []string{
			filepath.Join(tmp, "d1"),
			filepath.Join(tmp, "d2"),
			filepath.Join(tmp, "d3"),
		}},
		{opt: "-n13", glob: filepath.Join(tmp, "./*"), expects: []string{
			filepath.Join(tmp, "0"),
			filepath.Join(tmp, "1"),
			filepath.Join(tmp, "2"),
			filepath.Join(tmp, "3"),
			filepath.Join(tmp, "4"),
			filepath.Join(tmp, "5"),
			filepath.Join(tmp, "6"),
			filepath.Join(tmp, "7"),
			filepath.Join(tmp, "8"),
			filepath.Join(tmp, "9"),
			filepath.Join(tmp, "d1"),
			filepath.Join(tmp, "d2"),
			filepath.Join(tmp, "d3"),
		}},
		{opt: "-n1", glob: filepath.Join(tmp, "./?"), expects: []string{
			filepath.Join(tmp, "0"),
		}},
		{opt: "-n100", glob: filepath.Join(tmp, "./*/*"), expects: []string{
			filepath.Join(tmp, "d1", "f0"),
			filepath.Join(tmp, "d1", "f1"),
			filepath.Join(tmp, "d1", "f2"),
			filepath.Join(tmp, "d2", "f0"),
			filepath.Join(tmp, "d2", "f1"),
			filepath.Join(tmp, "d2", "f2"),
			filepath.Join(tmp, "d3", "f0"),
			filepath.Join(tmp, "d3", "f1"),
			filepath.Join(tmp, "d3", "f2"),
		}},
		{opt: "-n1", glob: filepath.Join(tmp, "1"), expects: []string{filepath.Join(tmp, "1")}},
		{glob: filepath.Join(tmp, "["), throw: true},
		{glob: filepath.Join(tmp, "x"), opt: "-n1", expects: []string{}},
	}

	for _, v := range data {
		ff.Options = []string{v.opt}
		actual, err := ff.GetPathes(v.glob)
		if v.throw {
			as.Error(err, v.glob)
			as.Empty(actual, v.glob)
		} else {
			as.NoError(err, v.glob)
			as.Equal(v.expects, actual, v.glob)
		}
	}

	ff.Command = "NOT_FOUND"
	ff.Options = []string{}
	t.Run("failed finder", func(t *testing.T) {
		a, e := ff.GetPathes(filepath.Join(tmp, "fail"))
		as.Nil(a)
		as.Error(e)
	})
}
