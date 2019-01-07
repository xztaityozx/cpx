type Config struct {
  FuzzyFinder FuzzyFinder
}

type FuzzyFinder struct {
  Command string
  Options []string
}
