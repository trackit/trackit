//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

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
