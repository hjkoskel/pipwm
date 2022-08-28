/*
TODO separate library later?
*/

package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/hjkoskel/gomonochromebitmap"
)

//Type that combines functions for SSD1306 I2C display
type SSD1306_i2c struct {
	address byte
	DispOn  bool
	hw      *os.File
	//TODO possible display buffer (for comparing prev vs next frame for more intelligent and faster display update)
}

//Draws fixed size 128x64 bitmap to buffer. TODO optimize number of call and build directly 0x40 data 0x40 data 0x40 data format
func BitmapToSSD1306Buffer(bitmap gomonochromebitmap.MonoBitmap, flipY bool) [1024]byte {
	result := [1024]byte{}
	if flipY {
		for x := 0; x < 128; x++ {
			for y := 0; y < 64; y++ {
				if bitmap.GetPix(127-x, 63-y) {
					result[x+128*int(y/8)] |= (1 << (byte(y) & 7))
				}
			}
		}
		return result
	}
	for x := 0; x < 128; x++ {
		for y := 0; y < 64; y++ {
			if bitmap.GetPix(x, y) {
				result[x+128*int(y/8)] |= (1 << (byte(y) & 7))
			}
		}
	}
	return result
}

func SSD1306BufferToBitmap(binArr [1024]byte, flipY bool) gomonochromebitmap.MonoBitmap {
	result := gomonochromebitmap.NewMonoBitmap(128, 64, false)
	for iv, v := range binArr {
		for bit := 0; bit < 8; bit++ {
			if 0 < v&(1<<bit) {
				result.SetPix(int(uint32(iv)&0x7F), bit+8*int(iv/128), true)
			}
		}
	}
	if flipY {
		result.FlipH()
		result.FlipV()
	}
	return result
}

//BW128x64Display is interface for real and "fake" display functionality
type BW128x64Display interface {
	DisplayOn() error
	DisplayOff() error
	ToggleOnOff() error
	FullDisplayUpdate(data [1024]byte) error
	Init() error
}

func InitSSD1306_i2c(deviceFile *os.File, address byte) (SSD1306_i2c, error) {
	result := SSD1306_i2c{hw: deviceFile, address: address}
	result.Init()
	return result, nil
}

func (p *SSD1306_i2c) DisplayOn() error {
	selErr := SelectI2CSlave(p.hw, p.address)
	if selErr != nil {
		return selErr
	}
	_, wErr := p.hw.Write([]byte{0, 0xaf}) // display on
	if wErr != nil {
		return fmt.Errorf("DisplayOn fail %v", wErr)
	}
	p.DispOn = true
	return nil
}

func (p *SSD1306_i2c) DisplayOff() error {
	selErr := SelectI2CSlave(p.hw, p.address)
	if selErr != nil {
		return selErr
	}
	_, wErr := p.hw.Write([]byte{0, 0xae}) // display off
	if wErr != nil {
		return fmt.Errorf("DisplayOff fail %v", wErr)
	}
	p.DispOn = false
	return nil
}

func (p *SSD1306_i2c) ToggleOnOff() error {
	if p.DispOn {
		return p.DisplayOff()
	}
	return p.DisplayOn()
}

//Updates with SSD1306 raw data
func (p *SSD1306_i2c) FullDisplayUpdate(data [1024]byte) error {
	//p.selectI2CSlave(p.display_i2caddr)
	errSelect := SelectI2CSlave(p.hw, p.address)
	if errSelect != nil {
		return errSelect
	}
	//
	seqList := [][]byte{
		[]byte{0x00, 0x21}, // 0x21 COMMAND SSD1306_COLUMNADDR
		[]byte{0x00, 0},    // Column start address
		[]byte{0x00, 127},  // Column end address

		[]byte{0x00, 0x22},         // 0x22 COMMAND SSD1306_PAGEADDR
		[]byte{0x00, 0},            // Start Page address
		[]byte{0x00, (64 / 8) - 1}, // End Page address
	}
	for _, seq := range seqList {
		_, errWrite := p.hw.Write(seq)
		if errWrite != nil {
			return errWrite
		}
	}

	/*
		p.hw.Write([]byte{0x00, 0x21}) // 0x21 COMMAND SSD1306_COLUMNADDR
		p.hw.Write([]byte{0x00, 0})    // Column start address
		p.hw.Write([]byte{0x00, 127})  // Column end address

		p.hw.Write([]byte{0x00, 0x22})         // 0x22 COMMAND SSD1306_PAGEADDR
		p.hw.Write([]byte{0x00, 0})            // Start Page address
		p.hw.Write([]byte{0x00, (64 / 8) - 1}) // End Page address
	*/

	for i := 0; i < 1024; i++ {
		_, errWrite := p.hw.Write([]byte{0x40, data[i]})
		if errWrite != nil {
			return errWrite
		}
	}
	return nil
}

func (p *SSD1306_i2c) Init() error {
	selErr := SelectI2CSlave(p.hw, p.address)
	if selErr != nil {
		return selErr
	}

	seqList := [][]byte{
		[]byte{0x00, 0xae}, // display off
		[]byte{0x00, 0xd5}, // clockdiv
		[]byte{0x00, 0x80},
		[]byte{0x00, 0xa8}, // multiplex
		[]byte{0x00, 0x3f},
		[]byte{0x00, 0xd3}, // offset
		[]byte{0x00, 0x00},
		[]byte{0x00, 0x40}, // startline
		[]byte{0x00, 0x8d}, // charge pump
		[]byte{0x00, 0x14},
		[]byte{0x00, 0x20}, // memory mode
		[]byte{0x00, 0x00},
		[]byte{0x00, 0xa1}, // segregmap
		[]byte{0x00, 0xc8}, // comscandec
		[]byte{0x00, 0xda}, // set com pins
		[]byte{0x00, 0x12},
		[]byte{0x00, 0x81}, // contrast
		[]byte{0x00, 0xcf},
		[]byte{0x00, 0xd9}, // precharge
		[]byte{0x00, 0xf1},
		[]byte{0x00, 0xdb}, // vcom detect
		[]byte{0x00, 0x40},
		[]byte{0x00, 0xa4}, // resume
		[]byte{0x00, 0xa6}, // normal (not inverted)
		[]byte{0x00, 0xaf}, // display on
	}
	for _, seq := range seqList {
		_, errWrite := p.hw.Write(seq)
		if errWrite != nil {
			return errWrite
		}
	}

	p.DispOn = true
	return nil
}

func SelectI2CSlave(f *os.File, address byte) error {
	//i2c_SLAVE := 0x0703
	_, _, errorcode := syscall.Syscall6(syscall.SYS_IOCTL, f.Fd(), 0x0703, uintptr(address), 0, 0, 0)
	if errorcode != 0 {
		return fmt.Errorf("Select I2C slave errcode %v", errorcode)
	}
	return nil
}
