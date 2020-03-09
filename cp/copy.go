package cp

import "io"

type Copy struct {
	src *Source
	dst *Destination
}

func New(src *Source, dst *Destination) Copy {
	return Copy{src: src, dst: dst}
}

// Run execute copy
func (cp Copy) Run() error {
	defer cp.src.Close()
	defer cp.dst.Close()

	_, err := io.Copy(cp.dst, cp.src)
	return err
}
