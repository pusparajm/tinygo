package main

import (
	"fmt"
	"flag"
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
	fmt.Fprintf(w, "	fs.SaveFile(\"%s\", [...]byte{", fn)
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
	fmt.Fprintf(w, "}[:])\n")
}

func main() {
	out := os.Stdout
	fmt.Fprintf(out, "package main\nimport \"tinygo.googlecode.com/hg/fs\"\n")
	fmt.Fprintf(out, "func init() {\n")
	for i := 0; i < flag.NArg(); i++ {
		doFile(out, flag.Arg(i))
	}
	fmt.Fprintf(out, "}\n")
}
