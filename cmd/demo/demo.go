package main

import (
	"image"
	"exp/draw"
//	"exp/draw/svga"
)

func main () {
	//ctxt, err := svga.NewScreen()
	//if ctxt == nil {
	//	log.Exitf("no screen: %v", err)
  //}
	//screen := ctxt.Screen()

	//img := screen.(*image.RGBA)
	img := image.NewRGBA(640,480)
	draw.Draw(img, image.Rect(10,10,100,100), image.White, image.ZP)

	cr, cg, cb, ca := image.White.RGBA()
	println("cr ", cr)
	println("cg ", cg)
	println("cb ", cb)
	println("ca ", ca)

	//ctxt.FlushImage()
	print("done.")
}
