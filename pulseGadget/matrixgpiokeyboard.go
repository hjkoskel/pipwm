/*
matrixgpiokeyboard.go

*/

package main

import (
	"fmt"
	"os"

	"github.com/hjkoskel/govattu"
	"github.com/hjkoskel/pipwm/pulsegenui"
	term "github.com/nsf/termbox-go"
)

type KeyboardInterface interface {
	Scan() SetOfPressedKeys
	Init()
}

type MatrixGpioKeyboard struct {
	DriveRows []uint8
	InputCols []uint8
	Buttons   [][]string //Coding: [driveRow][inputCol]
	prevState []string   //Compare here,
}

func (p *MatrixGpioKeyboard) Init() {
	for _, pin := range p.DriveRows {
		hw.PinMode(pin, govattu.ALTinput)
		hw.PullMode(pin, govattu.PULLdown)
	}
	for _, pin := range p.InputCols {
		hw.PinMode(pin, govattu.ALTinput)
		hw.PullMode(pin, govattu.PULLdown)
	}
}

type SetOfPressedKeys []string

func (p *SetOfPressedKeys) Equal(ref SetOfPressedKeys) bool {
	if len(*p) != len(ref) {
		return false //can not be same
	}
	for i, v := range *p {
		if ref[i] != v {
			return false
		}
	}
	return true
}

/*
Just scans situation now
*/
func (p *MatrixGpioKeyboard) Scan() SetOfPressedKeys {
	result := []string{}
	for driveRowNumber, drivePin := range p.DriveRows {
		hw.PinMode(drivePin, govattu.ALToutput)
		hw.PinSet(drivePin)

		keymask := hw.ReadAllPinLevels()

		names := p.Buttons[driveRowNumber]
		for inputColNumber, name := range names {
			if 0 < (keymask & (1 << p.InputCols[inputColNumber])) {
				result = append(result, name)
			}
		}
		hw.PinClear(drivePin)
		hw.PinMode(drivePin, govattu.ALTinput)
	}
	return result
}

/*****


*******/

type FakeMatrixKeyboard struct {
	/*DriveRows []uint8
	InputCols []uint8
	Buttons   [][]string //Coding: [driveRow][inputCol]
	prevState []string   //Compare here,
	*/

	//keyreader *bufio.Reader
	//scanner *bufio.Scanner
}

var keymapping map[term.Key]string

func (p *FakeMatrixKeyboard) Init() {
	//p.keyreader = bufio.NewReader(os.Stdin)
	//p.scanner = bufio.NewScanner(os.Stdin)

	/*
		keymapping = make(map[term.Key]string)
		keymapping[term.KeyEsc] =


		keymapping[57]=pulsegenui.CMDBTN_9
		keymapping[56]=pulsegenui.CMDBTN_8
		keymapping[55]=pulsegenui.CMDBTN_7
		keymapping[54]=pulsegenui.CMDBTN_6
		keymapping[53]=pulsegenui.CMDBTN_5
		keymapping[52]=pulsegenui.CMDBTN_4
		keymapping[51]=pulsegenui.CMDBTN_3
		keymapping[50]=pulsegenui.CMDBTN_2
		keymapping[49]=pulsegenui.CMDBTN_1
		keymapping[48]=pulsegenui.CMDBTN_0


		pulsegenui.CMDBTN_DECIMALPOINT, , pulsegenui.CMDBTN_ONOFF}}
	*/
	term.Init()
}

func (p *FakeMatrixKeyboard) Scan() SetOfPressedKeys {
	/*[]string{pulsegenui.CMDBTN_UP, pulsegenui.CMDBTN_9, pulsegenui.CMDBTN_8, pulsegenui.CMDBTN_7},
	[]string{pulsegenui.CMDBTN_DOWN, pulsegenui.CMDBTN_6, pulsegenui.CMDBTN_5, pulsegenui.CMDBTN_4},
	[]string{pulsegenui.CMDBTN_BACK, pulsegenui.CMDBTN_3, pulsegenui.CMDBTN_2, pulsegenui.CMDBTN_1},
	[]string{pulsegenui.CMDBTN_OK, pulsegenui.CMDBTN_DECIMALPOINT, pulsegenui.CMDBTN_0, pulsegenui.CMDBTN_ONOFF}}, //Coding: [driveRow][inputCol]
	*/

	/*
		ru, rulen, ruErr := p.keyreader.
			fmt.Printf("KEYBOARD: ru=%v rulen=%v ruErr=%v", ru, rulen, ruErr)
	*/

	/*
		var b = make([]byte, 1)
		os.Stdin.Read(b)
		fmt.Printf("b=%v\n", b)
	*/
	//fmt.Printf("SCANNED =%v", p.scanner.Scan())

	ev := term.PollEvent()
	fmt.Printf("Ev=%#v\n", ev)

	if ev.Type == term.EventKey {
		//ev.Key
		switch ev.Key {
		case term.KeyArrowUp:
			return SetOfPressedKeys{pulsegenui.CMDBTN_UP}
		case term.KeyArrowDown:
			return SetOfPressedKeys{pulsegenui.CMDBTN_DOWN}
		case term.KeyArrowLeft, term.KeyEsc:
			return SetOfPressedKeys{pulsegenui.CMDBTN_BACK}
		case term.KeyArrowRight, term.KeyEnter:
			return SetOfPressedKeys{pulsegenui.CMDBTN_OK}
		case term.KeyCtrlC:
			os.Exit(0)

		case term.KeySpace:
			return SetOfPressedKeys{pulsegenui.CMDBTN_ONOFF}

		}

		fmt.Printf("ch=%v\n", ev.Ch)
		switch ev.Ch {
		case 57:
			return SetOfPressedKeys{pulsegenui.CMDBTN_9}
		case 56:
			return SetOfPressedKeys{pulsegenui.CMDBTN_8}
		case 55:
			return SetOfPressedKeys{pulsegenui.CMDBTN_7}
		case 54:
			return SetOfPressedKeys{pulsegenui.CMDBTN_6}
		case 53:
			return SetOfPressedKeys{pulsegenui.CMDBTN_5}
		case 52:
			return SetOfPressedKeys{pulsegenui.CMDBTN_4}
		case 51:
			return SetOfPressedKeys{pulsegenui.CMDBTN_3}
		case 50:
			return SetOfPressedKeys{pulsegenui.CMDBTN_2}
		case 49:
			return SetOfPressedKeys{pulsegenui.CMDBTN_1}
		case 48:
			return SetOfPressedKeys{pulsegenui.CMDBTN_0}
		case 44, 46:
			return SetOfPressedKeys{pulsegenui.CMDBTN_DECIMALPOINT}
		}
	}

	return SetOfPressedKeys{}

}
