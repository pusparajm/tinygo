package main

import (
	"unsafe"
)

var crt *[25 * 80]uint16

type ColoredChar struct {
	back uint8
	fore uint8
	char uint8
}

const (
	black uint8 = iota
	blue
	green
	cyan
	red
	magenta
	brown
	white
)

func pokeScreen(c uint16, x, y int) {
	if crt == nil {
		// init on demand in case printf is called before
		// initialization runs.
		var mem uintptr = 0xb8000
		crt = (*[25 * 80]uint16)(unsafe.Pointer(mem))
		for i := range crt[0:] {
			crt[i] = 0
		}
	}

	if y > 24 || x > 79 {
		panic("writing off screen")
	}
	pos := y*80 + x
	if pos > cap(crt) {
		print("bad pos ", pos)
	}
	crt[pos] = c
}

type Screen struct {
	w, h int
	scr  []ColoredChar
}

func (s *Screen) set(x, y int, cc ColoredChar) {
	s.scr[s.w*y+x] = cc
}

func (s *Screen) flush() {
	var val uint16
	for i := 0; i < 80; i++ {
		for j := 0; j < 25; j++ {
			cc := s.scr[s.w*j+i]
			val = (uint16(cc.back)&0x0f)<<12 |
				(uint16(cc.fore)&0x0f)<<8 |
				uint16(cc.char)
			pokeScreen(val, i, j)
		}
	}
}

func newScreen(w, h int) *Screen {
	s := new(Screen)
	s.w = w
	s.h = h
	s.scr = make([]ColoredChar, w*h)
	return s
}

func intToAsc(i int) (c uint8) {
	if i >= 0 && i < 10 {
		c = 48 + uint8(i)
	} else if i >= 0 && i <= 15 {
		c = 97 + uint8(i) - 10
	} else {
		c = 32
	}
	return
}

func main() {
	screen := newScreen(80, 25)

	// show colors on the bottom
	for i := 0; i < 16; i++ {
		screen.set(i, 24, ColoredChar{uint8(i), black, ' '})
	}

	for i := 0; i < 16; i++ {
		screen.set(0, i+1, ColoredChar{black, white, intToAsc(i)})
		for j := 0; j < 16; j++ {
			screen.set(2+j, i+1, ColoredChar{black, white, uint8(i*16 + j)})
		}
	}
	screen.flush()
}
