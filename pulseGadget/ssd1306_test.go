package main

import (
	"fmt"
	"image"
	"testing"

	"github.com/hjkoskel/gomonochromebitmap"
	"github.com/stretchr/testify/assert"
)

func TestConvertPic(t *testing.T) {
	testpic := gomonochromebitmap.NewMonoBitmap(128, 64, false)
	testpic.Circle(image.Point{X: 17, Y: 32}, 16, true)
	testpic.Line(image.Point{X: 0, Y: 0}, image.Point{X: 127, Y: 63}, true)

	byteArrNoflip := BitmapToSSD1306Buffer(testpic, false)
	testpicNoFlip := SSD1306BufferToBitmap(byteArrNoflip, false)
	byteArrNoflipRef := BitmapToSSD1306Buffer(testpicNoFlip, false)

	assert.Equal(t, byteArrNoflip, byteArrNoflipRef)

	blockgraph := gomonochromebitmap.BlockGraphics{
		Clear: false, HaveBorder: true,
	}
	fmt.Printf("ORIGINAL\n%s\nCONVERSION\n%s\n\n", blockgraph.ToQuadBlockChars(&testpic), blockgraph.ToQuadBlockChars(&testpicNoFlip))

	byteArrFlip := BitmapToSSD1306Buffer(testpic, true)
	testpicFlip := SSD1306BufferToBitmap(byteArrFlip, true)
	byteArrFlipRef := BitmapToSSD1306Buffer(testpicFlip, true)

	fmt.Printf("FLIP\n%s\nCONVERSION\n%s\n\n", blockgraph.ToQuadBlockChars(&testpic), blockgraph.ToQuadBlockChars(&testpicNoFlip))

	assert.Equal(t, byteArrFlip, byteArrFlipRef)

}
