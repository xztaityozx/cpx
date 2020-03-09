package factory

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_Factory_Local(t *testing.T) {
	tmp := filepath.Join(os.TempDir(), "cpx")
	os.MkdirAll(tmp, 0755)
	defer os.RemoveAll(tmp)

	files := []string{}
	for i := 0; i < 10; i++ {
		files = append(files, filepath.Join(tmp, fmt.Sprint(i)))
	}

	for i, v := range files {
		ioutil.WriteFile(v, []byte(fmt.Sprintf("this is test file%d", i)), 0644)
	}

}
