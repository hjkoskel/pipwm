package main

import (
	"fmt"

	"github.com/hjkoskel/gomonochromebitmap"
)

type Fakedisplay128x64 struct {
}

func (p *Fakedisplay128x64) DisplayOn() error   { return nil }
func (p *Fakedisplay128x64) DisplayOff() error  { return nil }
func (p *Fakedisplay128x64) ToggleOnOff() error { return nil }
func (p *Fakedisplay128x64) FullDisplayUpdate(data [1024]byte) error {
	pic := SSD1306BufferToBitmap(data, false)

	blockgraph := gomonochromebitmap.BlockGraphics{
		Clear: false, HaveBorder: true,
	}
	fmt.Printf("%s\n", blockgraph.ToQuadBlockChars(&pic))
	return nil
}
func (p *Fakedisplay128x64) Init() error { return nil }
