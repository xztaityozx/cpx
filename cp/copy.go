package cp

import (
	"io"
	"os"
	"path/filepath"
)

type Copy struct {
	src *Source
	dst *Destination
}

func create(src *Source, dst *Destination) Copy {
	return Copy{src: src, dst: dst}
}

// Run execute copy
func (cp Copy) Run() error {
	defer cp.src.Close()
	defer cp.dst.Close()

	_, err := io.Copy(cp.dst, cp.src)
	return err
}

func Glob(sList, dList []string, recursive bool, prompt func(string) bool) ([]Copy, []string, error) {
	var rt []Copy
	var skip []string
	for _, src := range sList {
		for _, dst := range dList {
			err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				abs := filepath.Join(path)
				if filepath.Dir(abs) != src && !recursive {
					skip = append(skip, abs)
					return nil
				}

				// source is directory
				if info.IsDir() {
					return os.MkdirAll(filepath.Join(dst, path), info.Mode().Perm())
				}

				// source file is regular file
				if info.Mode().IsRegular() {
					s, err := File(abs)
					if err != nil {
						return err
					}
					// overwrite confirming
					d, err := Dst(filepath.Join(dst, path), prompt)
					if err != nil {
						return err
					}
					rt = append(rt, create(s, d))
				}

				skip = append(skip, abs)
				return nil
			})
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return rt, skip, nil
}
