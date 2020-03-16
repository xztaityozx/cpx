package cp

import (
	"io"
	"os"
	"path/filepath"

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

func Glob(glob string, recursive bool) ([]FileEntry, error) {
	var rt []FileEntry

	expanded, err := filepath.Glob(filepath.Clean(glob))
	if err != nil {
		return nil, err
	}

	visited := map[string]bool{}

	logrus.Info("-----------Start----------")

	for _, item := range expanded {
		logrus.Info("item: ", item)
		if visited[item] {
			logrus.Warn("continue...")
			continue
		}
		logrus.Info("filepath.Walk, item: ", item)
		err = filepath.Walk(item, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if visited[path] {
				logrus.Warn("continue...(filepath.Walk)")
				return nil
			}
			visited[item] = true

			d, n := filepath.Split(path)
			d = filepath.Clean(d)
			logrus.Info("d: ", d, " n: ", n)
			if d != item && !recursive {
				logrus.Warn("skipping...")
				return nil
			}
			logrus.Info("append to rt")
			rt = append(rt, FileEntry{base: d, path: n, info: info})

			return nil
		})
	}

	return rt, nil
}
