package main

import (
	"fmt"
	"log"
	"io"
	"os"
)

func doFile(w io.Writer, fn string) {
	f, err := os.Open(fn, os.O_RDONLY, 0)
	if err != nil {
		log.Print("Could not read ", fn, ": ", err)
		return
	}
	fmt.Fprintf(w, "var fsString string = \"")
	for {
		var buf [512]byte
		n, err := f.Read(buf[:])
		if err != nil {
			if err != os.EOF {
				log.Print("Error reading file: ", err)
			}
			break
		}
		for i := 0; i < n; i++ {
			if isPrintable(buf[i]) {
				fmt.Fprintf(w, "%c", buf[i])
			} else {
				fmt.Fprintf(w, "0x%x", buf[i])
			}
		}
	}
	fmt.Fprintf(w, "\"\n")
}

func isPrintable(x byte) bool {
	if x >= 'a' && x <= 'z' {
		return true
	}
	if x >= 'A' && x <= 'Z' {
		return true
	}
	if x >= '0' && x <= '9' {
		return true
	}
	return false
}

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Should have only one argument.\n")
		os.Exit(1)
	}

	out := os.Stdout
	fmt.Fprintf(out, "// Generated file: do not edit. Use mkfs instead.\n")
	fmt.Fprintf(out, "package main\n")
	fmt.Fprintf(out, `
import (
//	"reflect"
//	"unsafe"
	"archive/zip"
	"os"
)
`)

	doFile(out, os.Args[1])

	fmt.Fprintf(out, `
var fs *zip.Reader

//type sliceReaderAt []byte
//
//func (r sliceReaderAt) ReadAt(b []byte, off int64) (int, os.Error) {
//  copy(b, r[int(off):int(off)+len(b)])
//  return len(b), nil
//}

type stringReaderAt string
func (r stringReaderAt) ReadAt(b []byte, off int64) (int, os.Error) {
	copy(b, r[int(off):int(off)+len(b)])
	return len(b), nil
}

func init() {
//	sx := (*reflect.StringHeader)(unsafe.Pointer(&fsString))
//	var x [0]byte
//	b := x[:]
//	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
//	bx.Data = sx.Data
//	bx.Len = len(fsString)
//	bx.Cap = len(fsString)
	var err os.Error
	fs, err = zip.NewReader(stringReaderAt(fsString), int64(len(fsString)))
	if err != nil {
		panic("could not read a zipfile from the filesystem string")
	}
}
`)
}
