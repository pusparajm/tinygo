package main

import "tinygo.googlecode.com/hg/fs"

func initFS() {
	fs.SaveFile("test", []byte("test\n"))
}
