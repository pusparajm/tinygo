package main

import (
	"fmt"
	"flag"
	"log"
	"io"
	"os"
	"unicode"
	"strings"
)

func doFile(w io.Writer, fn string) {
	f, err := os.Open(fn, os.O_RDONLY, 0)
	if err != nil {
		log.Print("Could not read ", fn, ": ", err)
		return
	}
	fmt.Fprintf(w, "const unsigned char data_%s[] = {", toC(fn))
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
			fmt.Fprintf(w, "0x%x,", buf[i])
		}
	}
	fmt.Fprintf(w, "};\n")
}

func toC(s string) string {
	return strings.Map(
		func (ch int) int {
			if (unicode.IsDigit(ch) || unicode.IsLetter(ch)) {
				return ch
			}
			return '_'
		}, s)
}

func main() {
	out := os.Stdout
	fmt.Fprintf(out, "// Generated file: do not edit. Use mkfs instead.\n")
	fmt.Fprintf(out, "package fs\n")
	fmt.Fprintf(out, "/*\n")
	for i := 0; i < flag.NArg(); i++ {
		doFile(out, flag.Arg(i))
	}
	fmt.Fprintf(out, "*/\n")
	fmt.Fprintf(out, "import \"C\"\n")
	fmt.Fprintf(out, `
import (
	"reflect"
	"unsafe"
)
var FileMap map[string][]byte

func fix(in []C.uchar) []byte {
	inx := (*reflect.SliceHeader)(unsafe.Pointer(&in))
	var buf [1]byte
	s := buf[:]
	sx := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	sx.Data = inx.Data
	sx.Len = inx.Len
	sx.Cap = inx.Cap
	return s
}
`)

	fmt.Fprintf(out, "func init() {\n")
	fmt.Fprintf(out, "	FileMap = make(map[string][]byte)\n")
	for i := 0; i < flag.NArg(); i++ {
		file := flag.Arg(i)
		fmt.Fprintf(out, "	FileMap[\"%s\"] = fix(C.data_%s[:])\n", file, toC(file))
	}
	fmt.Fprintf(out, "}\n")
}
