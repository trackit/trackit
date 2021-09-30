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

// Package util implements a few generic utilities that are widely useful for the entire server code
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
