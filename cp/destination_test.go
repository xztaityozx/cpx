package cp

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dst(t *testing.T) {
	yes := func(s string) bool { return true }
	no := func(s string) bool { return false }
	as := assert.New(t)
	tmp := filepath.Join(os.TempDir(), "cpx")
	os.MkdirAll(tmp, 0755)
	defer os.RemoveAll(tmp)

	t.Run("regular file(not exist)", func(t *testing.T) {
		rf := filepath.Join(tmp, "rf")
		d, err := Dst(rf, yes, false)
		as.NotNil(d)
		as.NoError(err)
		as.FileExists(rf)
		os.Remove(rf)
	})

	t.Run("regular file", func(t *testing.T) {
		rf := filepath.Join(tmp, "rf")
		ioutil.WriteFile(rf, []byte("file"), 0644)
		d, err := Dst(rf, yes, false)
		as.NotNil(d)
		as.NoError(err)
		as.FileExists(rf)
		os.Remove(rf)
	})

	t.Run("not regular file", func(t *testing.T) {
		d, err := Dst(os.DevNull, yes, false)
		as.Nil(d)
		as.Error(err)
	})

	t.Run("is directory", func(t *testing.T) {
		d, err := Dst(tmp, yes, false)
		as.Nil(d)
		as.Error(err)
	})

	t.Run("force", func(t *testing.T) {
		rf := filepath.Join(tmp, "rf")
		ioutil.WriteFile(rf, []byte("file"), 0644)
		d, err := Dst(rf, yes, true)
		as.NotNil(d)
		as.NoError(err)
		as.FileExists(rf)
	})

	t.Run("cancel", func(t *testing.T) {
		rf := filepath.Join(tmp, "rf")
		ioutil.WriteFile(rf, []byte("file"), 0644)
		d, err := Dst(rf, no, false)
		as.Nil(d)
		as.Error(err)
		os.Remove(rf)
	})

}
