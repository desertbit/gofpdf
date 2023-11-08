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
	Close() error
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

func (m *memBuffer) Close() error {
	return nil
}

type fileBuffer struct {
	f *os.File
}

func newTempFileBuffer(dir string) (fb *fileBuffer, err error) {
	f, err := os.CreateTemp(dir, "")
	if err != nil {
		return
	}

	return &fileBuffer{f: f}, nil
}

func newFileBuffer(path string) (fb *fileBuffer, err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}

	return &fileBuffer{f: f}, nil
}

func (f *fileBuffer) Read(p []byte) (n int, err error) {
	f.f.Seek(0, io.SeekStart)

	return f.f.Read(p)
}

func (f *fileBuffer) Write(p []byte) (n int, err error) {
	f.f.Seek(0, io.SeekEnd)

	return f.f.Write(p)
}

func (f *fileBuffer) WriteString(s string) (n int, err error) {
	f.f.Seek(0, io.SeekEnd)

	return f.f.WriteString(s)
}

func (f *fileBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	f.f.Seek(0, io.SeekEnd)

	return f.f.ReadFrom(r)
}

func (f *fileBuffer) Truncate(size int64) error {
	return f.f.Truncate(size)
}

func (f *fileBuffer) String() string {
	return string(f.Bytes())
}

func (f *fileBuffer) Bytes() []byte {
	f.f.Seek(0, io.SeekStart)

	data, err := io.ReadAll(f.f)
	if err != nil {
		panic(err)
	}

	return data
}

func (f *fileBuffer) Len() int {
	info, err := f.f.Stat()
	if err != nil {
		panic(err)
	}

	return int(info.Size())
}

func (f *fileBuffer) Printf(fmtStr string, args ...interface{}) {
	f.f.WriteString(fmt.Sprintf(fmtStr, args...))
}

func (f *fileBuffer) WriteTo(w io.Writer) (n int64, err error) {
	f.f.Seek(0, io.SeekEnd)

	return io.Copy(w, f.f)
}

func (f *fileBuffer) Close() error {
	return f.f.Close()
}
