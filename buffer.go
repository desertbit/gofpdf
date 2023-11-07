// Copyright 2023 skaldesh
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gofpdf

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var (
	_ buffer = &memBuffer{}
	_ buffer = &fileBuffer{}
)

type buffer interface {
	io.ReadWriter
	io.StringWriter
	io.ReaderFrom
	io.WriterTo

	Truncate(size int64) error
	String() string
	Bytes() []byte
	Len() int

	Printf(fmtStr string, args ...interface{})
}

type memBuffer struct {
	*bytes.Buffer
}

func newMemBuffer() *memBuffer {
	return &memBuffer{Buffer: &bytes.Buffer{}}
}

func (m *memBuffer) Truncate(size int64) error {
	m.Buffer.Truncate(int(size))
	return nil
}

func (m *memBuffer) Printf(fmtStr string, args ...interface{}) {
	m.Buffer.WriteString(fmt.Sprintf(fmtStr, args...))
}

type fileBuffer struct {
	*os.File
}

func newTempFileBuffer(dir string) (fb *fileBuffer, err error) {
	f, err := os.CreateTemp(dir, "")
	if err != nil {
		return
	}

	return &fileBuffer{File: f}, nil
}

func newFileBuffer(path string) (fb *fileBuffer, err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}

	return &fileBuffer{File: f}, nil
}

func (f *fileBuffer) String() string {
	return string(f.Bytes())
}

func (f *fileBuffer) Bytes() []byte {
	data, err := io.ReadAll(f.File)
	if err != nil {
		panic(err)
	}

	return data
}

func (f *fileBuffer) Len() int {
	info, err := f.File.Stat()
	if err != nil {
		panic(err)
	}

	return int(info.Size())
}

func (f *fileBuffer) Printf(fmtStr string, args ...interface{}) {
	f.File.WriteString(fmt.Sprintf(fmtStr, args...))
}

func (f *fileBuffer) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, f.File)
}
