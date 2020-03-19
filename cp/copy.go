package cp

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
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

type FileEntry struct {
	base string
	path string
	info os.FileInfo
}

type PromptFunc func(string) bool

func getGlobRoot(glob string) (string, error) {
	path := glob
	for _, err := os.Stat(path); err != nil; {
		path = filepath.Dir(path)
		logrus.Info(path)
	}
	return path, nil
}

func Glob(glob string, recursive bool) ([]FileEntry, error) {
	var rt []FileEntry

	info, err := os.Stat(glob)
	base := glob
	for _, err := os.Stat(base); err != nil; base = filepath.Dir(base) {
	}
	visited := map[string]bool{}
	walk := func(item string) error {
		return filepath.Walk(item, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			path = filepath.Clean(path)
			if visited[path] {
				return nil
			}
			visited[path] = true

			rt = append(rt, FileEntry{base: base, path: strings.TrimPrefix(path, base), info: info})

			return nil
		})
	}
	if err != nil {
		// is glob?
		matches, err := filepath.Glob(glob)
		if err != nil {
			// globでもpathでもない
			return nil, err
		}

		for _, item := range matches {
			// globにマッチしたもの
			item = filepath.Clean(item)
			if visited[item] {
				continue
			}
			visited[item] = true

			if recursive {
				if err := walk(item); err != nil {
					return nil, err
				}
			} else {
				if info, err := os.Stat(item); err == nil && !info.IsDir() {
					rt = append(rt, FileEntry{base: base, path: strings.TrimPrefix(item, base), info: info})
				}
			}
		}
	} else {
		if info.IsDir() {
			if recursive {
				if err := walk(glob); err != nil {
					return nil, err
				}
			} else {
				entries, err := ioutil.ReadDir(glob)
				if err != nil {
					return nil, err
				}
				for _, entry := range entries {
					if !entry.IsDir() {
						rt = append(rt, FileEntry{
							base: filepath.Dir(glob),
							path: filepath.Dir(filepath.Join(glob, entry.Name())),
							info: entry})
					}
				}
			}
		} else {
			rt = append(rt, FileEntry{base: glob, path: info.Name(), info: info})
		}
	}

	return rt, nil
}
