/*
Pulse gadget.
Simple gadget using  128x64 oled display and 4x4 matrix keyboad
*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hjkoskel/goDisKeySystem/hwLayer"
	"github.com/hjkoskel/gomonochromebitmap"
	"github.com/hjkoskel/govattu"

	"github.com/hjkoskel/pipwm/pulseGenUi"

)

const (
	I2CADDRESS_OLED = 0x3c
)
const (
	I2CDEVICEFILE = "/dev/i2c-1"
)

const (
	KEYPOLLINTERVALMS=50
)

/*
Nice keyboard helper
*/

type MatrixGpioKeyboard struct {
	DriveRows []uint8
	InputCols []uint8
	Buttons   [][]string //Coding: [driveRow][inputCol]
	prevState []string   //Compare here,
}

func (p *MatrixGpioKeyboard) Init() {
	for _, pin := range p.DriveRows {
		govattu.PinMode(pin, govattu.ALTinput)
		govattu.PullMode(pin, govattu.PULLdown)
	}
	for _, pin := range p.InputCols {
		govattu.PinMode(pin, govattu.ALTinput)
		govattu.PullMode(pin, govattu.PULLdown)
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
		govattu.PinMode(drivePin, govattu.ALToutput)
		govattu.PinSet(drivePin)

		//time.Sleep(time.Millisecond * 200)
		keymask := govattu.ReadAllPinLevels()
		//fmt.Printf("Drive:%v  regs:%b\n", driveRowNumber, keymask)

		names := p.Buttons[driveRowNumber]
		for inputColNumber, name := range names {
			if 0 < (keymask & (1 << p.InputCols[inputColNumber])) {
				result = append(result, name)
			}
		}
		govattu.PinClear(drivePin)
		govattu.PinMode(drivePin, govattu.ALTinput)
	}
	//fmt.Printf("Keys hit %#v\n", result)
	return result
}

/*
gives channels
*/

/*
func runWithRealHardware(bitmapCh chan gomonochromebitmap.MonoBitmap, keysCh chan string) (govattu.RaspiHw, error) {
	return hw, nil
}
*/

func main() {
	BitmapCh := make(chan gomonochromebitmap.MonoBitmap, 1)
	CmdCh := make(chan string, 3) //Key pressesses. This gadget takes only one press per time
	ui := pulseGenUi.PulsGenUi{Bitmap: BitmapCh, Cmd: CmdCh, Simulate: false}

	hw, err := govattu.Open()
	if err != nil {
		fmt.Printf("Raspberry hardware fail %v\n",err.Error())
		os.Exit(-1)
	}
	defer hw.Close()

	hardwareKeyboard := MatrixGpioKeyboard{
		//	BCM codes
		DriveRows: []uint8{12, 16, 20, 21},
		InputCols: []uint8{26, 19, 13, 6},
		Buttons: [][]string{
			[]string{pulseGenUi.CMDBTN_UP, pulseGenUi.CMDBTN_9, pulseGenUi.CMDBTN_8, pulseGenUi.CMDBTN_7},
			[]string{pulseGenUi.CMDBTN_DOWN, pulseGenUi.CMDBTN_6, pulseGenUi.CMDBTN_5, pulseGenUi.CMDBTN_4},
			[]string{pulseGenUi.CMDBTN_BACK, pulseGenUi.CMDBTN_3, pulseGenUi.CMDBTN_2, pulseGenUi.CMDBTN_1},
			[]string{pulseGenUi.CMDBTN_OK, pulseGenUi.CMDBTN_DECIMALPOINT, pulseGenUi.CMDBTN_0, pulseGenUi.CMDBTN_ONOFF}}, //Coding: [driveRow][inputCol]
	}

	hardwareKeyboard.Init()
	go func() {
		prevKeyStatus := hardwareKeyboard.Scan()
		for {
			keyStatus := hardwareKeyboard.Scan()
			time.Sleep(KEYPOLLINTERVALMS * time.Millisecond)
			if !keyStatus.Equal(prevKeyStatus) {
				prevKeyStatus = keyStatus
				if len(keyStatus) == 0 {
					CmdCh <- pulseGenUi.CMDBTN_RELEASE //Changed and nothing pressed. report as released
				} else {
					for _, cmd := range keyStatus {
						CmdCh <- cmd
					}
				}
			}
		}
	}()
	//Display init
	fDisplay, errI2CHardware := os.OpenFile(I2CDEVICEFILE, os.O_RDWR, 0600)
	if errI2CHardware != nil {
		fmt.Printf("I2C Hardware open error %v", errI2CHardware.Error())
		os.Exit(-1)
	}

	oled, errOled := hwLayer.InitSSD1306_i2c(fDisplay, I2CADDRESS_OLED)
	if errOled != nil {
		fmt.Printf("I2C display init error %v", errOled.Error())
		os.Exit(-1)
	}

	go func() {
		for {
			oled.FullDisplayUpdate(hwLayer.BitmapToSSD1306Buffer(<-BitmapCh, false))
		}
	}()

	ui.Run() //This is where "business logic is ticking".

}
