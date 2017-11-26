package util

import (
	"errors"
)

type FixedBuffer []byte

var (
	ErrUnderflow = errors.New("underflow writing to FixedBuffer")
	ErrOverflow  = errors.New("overflow writing to FixedBuffer")
)

func (fb FixedBuffer) WriteAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, ErrUnderflow
	} else if a := int64(len(fb)); off > a {
		return 0, ErrOverflow
	} else if l := int64(len(p)); off+l > a {
		n, _ := fb.WriteAt(p[:off+l-a], off)
		return n, ErrOverflow
	} else {
		for i := range p {
			fb[off+int64(i)] = p[i]
		}
		return int(l), nil
	}
}
