package fs

import (
	"os"
	"io"
)

var theFS map[string]*file
var ro = false

type file struct {
	data []byte
}

func init() {
	theFS = make(map[string]*file)
}

func SaveFile(name string, data []byte) os.Error {
	if !ro {
		theFS[name] = &file{ data: data }
	} else {
		return os.NewError("filestore is read only")
	}
	return nil
}

func SetRO() {
	ro = true
}

type fileContext struct {
	pos int
	file *file
}

func Open(file string) (io.ReadCloser, os.Error) {
	f, ok := theFS[file]
	if ! ok {
		return nil, os.NewError("file not found")
	}
	return &fileContext{ file: f }, nil
}

// implements io.Reader
func (f *fileContext)Read(p []byte) (n int, err os.Error) {
	avail := len(f.file.data) - f.pos
	toread := len(p)

	if toread == 0 {
		return 0, os.EOF
	}
	if toread > avail {
		toread = avail
	}

	copy(p, f.file.data[f.pos:f.pos+toread])
	f.pos += toread
	return toread, nil
}

// implements io.Closer
func (f *fileContext)Close() os.Error {
	f.pos = 0
	return nil
}
