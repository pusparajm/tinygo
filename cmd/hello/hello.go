// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"archive/zip"
	"./fs/_obj/fs"
)

type readable []byte

func NewReadable(buf []byte) readable {
	return buf
}

func (r readable)ReadAt(p []byte, off int64) (n int, err os.Error) {
	o := int(off)
	copy(p, r[o:o+len(p)])
	return len(p), nil
}

func main() {
	println("Hello world from tiny!")
	println()

	a := 3.14159
	println("Floating point is initialized:")
	println("pi: ", a)
	println("2*pi: ", 2*a)
	println()

	println("Channels work:")
	Sieve()

	println("Read from fs:")
	b := fs.FileMap["/tmp/test.zip"]
	n, err := zip.NewReader(NewReadable(b), int64(len(b)))
	if err != nil {
		panic(err)
	}
	for _, x := range n.File {
		println("File:", x.Name)
	}
}

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in // Receive value of new variable 'i' from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to channel 'out'.
		}
	}
}

// The prime sieve: Daisy-chain Filter processes together.
func Sieve() {
	ch := make(chan int) // Create a new channel.
	go Generate(ch)      // Start Generate() as a subprocess.
	for i := 0; i < 10; i++ {
		prime := <-ch
		print("Prime ", i, ": ", prime, "\n")
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
}
