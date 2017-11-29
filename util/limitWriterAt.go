package util

import (
	"errors"
	"io"
)

var (
	ErrInvalidSize  = errors.New("Creating LimitWriterAt with invalid size.")
	ErrWriteLimited = errors.New("Write was limited by LimitWriterAt.")
)

type LimitWriterAt struct {
	writer io.WriterAt
	size   int64
}

func NewLimitWriterAt(w io.WriterAt, s int64) *LimitWriterAt {
	if s < 0 {
		panic(ErrInvalidSize)
	}
	return &LimitWriterAt{
		writer: w,
		size:   s,
	}
}

func (lwa LimitWriterAt) WriteAt(p []byte, off int64) (int, error) {
	l := int64(len(p))
	over := (off + l) - lwa.size
	if over > 0 {
		limit := l - over
		if n, err := lwa.writer.WriteAt(p[:limit], off); err != nil {
			return n, err
		} else {
			return n, ErrWriteLimited
		}
	} else {
		return lwa.writer.WriteAt(p, off)
	}
}
