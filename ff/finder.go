package ff

import (
	"fmt"

	"github.com/b4b4r07/go-finder"
	"github.com/b4b4r07/go-finder/source"
)

type FuzzyFinder struct {
	// /path/to/fuzzy-finder-command
	Command string
	// options for fuzzy-finder
	Options []string
}

// Return yes/no
func (ff FuzzyFinder) YesNo(prompt string) bool {
	const yes = "yes"
	fmt.Println(prompt)
	f, err := finder.New(append([]string{ff.Command}, ff.Options...)...)
	if err != nil {
		return false
	}

	f.Read(source.Slice([]string{yes, "no"}))
	if res, err := f.Run(); err != nil || len(res) != 1 || res[0] != yes {
		return false
	}
	return true
}

//func (ff FuzzyFinder) GetPathes() ([]string, error) {
//}
