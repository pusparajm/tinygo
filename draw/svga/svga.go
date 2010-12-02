// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package implements an SVGA backend for the exp/draw package.
//
// It is not working yet.

package svga

import (
	"syscall"
	"exp/draw"
	"image"
	"os"
	"unsafe"
)

const (
	windowHeight = 768
	windowWidth  = 1024
)

type screen struct {
	img	*image.RGBA
	flush	chan bool
	eventc	chan interface {}
}

func (s *screen) Screen() draw.Image { return s.img }

func (s *screen) FlushImage() {
	// We do the send (the <- operator) in an expression context, rather than in
	// a statement context, so that it does not block, and fails if the buffered
	// channel is full (in which case there already is a flush request pending).
	_ = s.flush <- false
}

func (s *screen) Close() os.Error {
	// tell the flusher to clean up and exit
	close(s.flush)
	return nil
}

// http://wiki.osdev.org/BGA

const (
	VBE_DISPI_INDEX_ID uint16 = iota
	VBE_DISPI_INDEX_XRES
	VBE_DISPI_INDEX_YRES
	VBE_DISPI_INDEX_BPP
	VBE_DISPI_INDEX_ENABLE
	VBE_DISPI_INDEX_BANK
	VBE_DISPI_INDEX_VIRT_WIDTH
	VBE_DISPI_INDEX_VIRT_HEIGHT
	VBE_DISPI_INDEX_X_OFFSET
	VBE_DISPI_INDEX_Y_OFFSET
	VBE_DISPI_IOPORT_INDEX uint16 = 0x01ce
	VBE_DISPI_IOPORT_DATA uint16 = 0x01cf
	VBE_DISPI_ID4 uint16 = 0xB0C4
	VBE_DISPI_LFB_ENABLED uint16 = 0x40
	VBE_DISPI_NOCLEARMEM uint16 = 0x80
	CGA_PORT uint16 = 0x3d8
	VBE_DISPI_LFB_PHYSICAL_ADDRESS = 0xE0000000
)

const (
	VBE_DISPI_DISABLED uint16 = iota
	VBE_DISPI_ENABLED
)

func bgaWriteRegister(index, word uint16) {
	syscall.Outw(VBE_DISPI_IOPORT_INDEX, index)
	syscall.Outw(VBE_DISPI_IOPORT_DATA, word)
}

func bgaReadRegister(index uint16) uint16 {
	syscall.Outw(VBE_DISPI_IOPORT_INDEX, index)
	return syscall.Inw(VBE_DISPI_IOPORT_DATA)
}

func bgaIsAvailable() bool {
    return (bgaReadRegister(VBE_DISPI_INDEX_ID) == VBE_DISPI_ID4)
}

func bgaSetVideoMode(Width, Height, BitDepth int, UseLinearFrameBuffer, ClearVideoMemory bool) {
	bgaWriteRegister(VBE_DISPI_INDEX_ENABLE, VBE_DISPI_DISABLED)
	bgaWriteRegister(VBE_DISPI_INDEX_XRES, uint16(Width))
	bgaWriteRegister(VBE_DISPI_INDEX_YRES, uint16(Height))
	bgaWriteRegister(VBE_DISPI_INDEX_BPP, uint16(BitDepth))

	var flags uint16 = VBE_DISPI_ENABLED
	if (UseLinearFrameBuffer) {
		flags |= VBE_DISPI_LFB_ENABLED
	}
	if (! ClearVideoMemory) {
		flags |= VBE_DISPI_NOCLEARMEM
	}
	bgaWriteRegister(VBE_DISPI_INDEX_ENABLE, flags)
}
 
func bgaSetBank(bank uint16) {
    bgaWriteRegister(VBE_DISPI_INDEX_BANK, bank);
}

var vgamem *[windowWidth * windowHeight * 3]uint8

func (s *screen) flusher() {
	// init VGA mode
	if (!bgaIsAvailable()) {
		panic("No Bochs-compatible VGA adapter found")
	} else {
		println("Bochs VGA is available.")
	}

	var savedCGA byte
  if false {
	savedCGA = syscall.Inb(CGA_PORT)
	var VBE_DISPI_BPP_24 = 0x18
	bgaSetVideoMode(windowWidth, windowHeight, VBE_DISPI_BPP_24, true, true)
	}

	if false {
	println("Getting point to the screen")
	var base uintptr = VBE_DISPI_LFB_PHYSICAL_ADDRESS
	vgamem = (*[windowWidth * windowHeight * 3]uint8)(unsafe.Pointer(base))

	for _ = range s.flush {
		i := 0
		for i < cap(vgamem) {
			vgamem[i] = 0xff
			vgamem[i+1] = 0xff
			vgamem[i+2] = 0xff
			i += 3
		}
		println("Finished a flush")
	}
	}
	println("Finished flusher goroutine")

	// go back to text mode
	syscall.Outb(CGA_PORT, savedCGA)
}

// eventually we'll have something writing on this channel...
func (s *screen) EventChan() <-chan interface{} { return s.eventc }

// NewScreen returns a new draw.Window, backed by the SVGA display
func NewScreen() (draw.Window, os.Error) {
	s := new(screen)
	s.img = image.NewRGBA(windowWidth, windowHeight)
	s.eventc = make(chan interface{}, 16)
	s.flush = make(chan bool, 1)
	go s.flusher()
	return s, nil
}
