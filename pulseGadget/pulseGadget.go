/*
Pulse gadget.
Simple gadget using  128x64 oled display and 4x4 matrix keyboad


TODO ELOKUUSSA
- pulsegenui + tekstimoodi 128x64. tai yksittäiset komennot
	- Sekä sim että ei simulaatio (interfacet)
- Juurihakemistossa vain main funktio ja komentoriviflägien parsinta?
- graafinen SDL gadgetti alihakemistoon

EIKU TEKEE JUUREEN SELLAISEN ETTÄ SILLÄ SAA PYÖRIMÄÄN.
pulsegenui:hin niitä mitä tarvii graafiseen käyttöliittymäänkin
	- simulaatio vain?
	- raspin työpöydällä?

- SIIRRÄ HARDWARE PWM jutut govattuun ja pulseratioiden laskeminen yms!!!

*/

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hjkoskel/gomonochromebitmap"
	"github.com/hjkoskel/govattu"
	"github.com/hjkoskel/pipwm/pulsegenui"
)

const (
	I2CADDRESS_OLED = 0x3c
)
const (
	I2CDEVICEFILE = "/dev/i2c-1"
)

const (
	KEYPOLLINTERVALMS = 50
)

/*
gives channels
*/
func setPwm0Generator(cmd pulsegenui.PulseHardwareCommand) error {
	if !(0 < cmd.Rf.On) || !cmd.OutputEnabled { //SHUTDOWN TO LO
		hw.SetPWM0Lo()
		return nil
	}

	if !(0 < cmd.Rf.Off) { //SHUTDOWN TO HI
		hw.SetPWM0Hi()
		return nil
	}

	//pulseTiming := govattu.RfPulseSettings{On: cmd.Rf.On, Off: cmd.Rf.Off}
	pulseSettings, pulseErr := cmd.Rf.GetSettings()
	if pulseErr != nil {
		return pulseErr
	}

	if pulseSettings.Pwmr <= pulseSettings.Pwm {
		return fmt.Errorf("Cant do requested on/off  %v/%v us pulse", float64(cmd.Rf.On.Nanoseconds())/1000.0, float64(cmd.Rf.Off.Nanoseconds())/1000.0)
	}

	if 4095 < pulseSettings.Pwmc {
		return fmt.Errorf("Cant do requested on/off  %v/%v us pulse pwmc can not be divided more", float64(cmd.Rf.On.Nanoseconds())/1000.0, float64(cmd.Rf.Off.Nanoseconds())/1000.0)
	}

	hw.SetToHwPWM0(&pulseSettings)
	return nil
}

/*
func runWithRealHardware(bitmapCh chan gomonochromebitmap.MonoBitmap, keysCh chan string) (govattu.RaspiHw, error) {
	return hw, nil
}
*/

var hw govattu.Vattu

func main() {

	pSimHw := flag.Bool("simhw", false, "simulate ALL hardware")
	pSimKeyboard := flag.Bool("simkey", false, "simulate 4x4 matrix keyboard with console keypresses")
	pSimDisplay := flag.Bool("simdisp", false, "simulate ssd1306 with console printout")
	flag.Parse()

	//fmt.Printf("pSimHw=%v,pSimKeyboard=%v,pSimDisplay=%v\n", *pSimHw, *pSimKeyboard, *pSimDisplay)
	//return

	if *pSimHw {
		fmt.Printf("SIMULATING ALL HARWARE\n")
	} else {
		if *pSimDisplay {
			fmt.Printf("SIMULATING DISPLAY\n")
		}
		if *pSimKeyboard {
			fmt.Printf("SIMULATING KEYBOARD\n")
		}
	}

	BitmapCh := make(chan gomonochromebitmap.MonoBitmap, 1)
	CmdCh := make(chan string, 3) //Key pressesses. This gadget takes only one press per time
	ui := pulsegenui.PulsGenUi{Bitmap: BitmapCh, Cmd: CmdCh, RfCmd: make(chan pulsegenui.PulseHardwareCommand, 10)}

	var errhw error

	if *pSimHw {
		hw = &govattu.DoNothingPi{}
	} else {
		hw, errhw = govattu.Open()
		if errhw != nil {
			fmt.Printf("Raspberry hardware fail %v\n", errhw.Error())
			os.Exit(-1)
		}
	}
	defer hw.Close()

	var hardwareKeyboard KeyboardInterface

	if *pSimHw || *pSimKeyboard {
		hardwareKeyboard = KeyboardInterface(&FakeMatrixKeyboard{})
	} else {
		hardwareKeyboard = KeyboardInterface(&MatrixGpioKeyboard{
			//	BCM codes
			DriveRows: []uint8{12, 16, 20, 21},
			InputCols: []uint8{26, 19, 13, 6},
			Buttons: [][]string{
				[]string{pulsegenui.CMDBTN_UP, pulsegenui.CMDBTN_9, pulsegenui.CMDBTN_8, pulsegenui.CMDBTN_7},
				[]string{pulsegenui.CMDBTN_DOWN, pulsegenui.CMDBTN_6, pulsegenui.CMDBTN_5, pulsegenui.CMDBTN_4},
				[]string{pulsegenui.CMDBTN_BACK, pulsegenui.CMDBTN_3, pulsegenui.CMDBTN_2, pulsegenui.CMDBTN_1},
				[]string{pulsegenui.CMDBTN_OK, pulsegenui.CMDBTN_DECIMALPOINT, pulsegenui.CMDBTN_0, pulsegenui.CMDBTN_ONOFF}}, //Coding: [driveRow][inputCol]
		})
	}

	hardwareKeyboard.Init()

	if *pSimHw || *pSimKeyboard {
		go func() {
			for {
				keyStatus := hardwareKeyboard.Scan()
				for _, cmd := range keyStatus {
					CmdCh <- cmd
				}
				//time.Sleep(time.Millisecond * 100)
			}
		}()

		go func() {
			for {
				if len(CmdCh) == 0 {
					CmdCh <- pulsegenui.CMDBTN_RELEASE
				}
				//time.Sleep(KEYPOLLINTERVALMS * time.Millisecond * 10)
				time.Sleep(time.Millisecond * 100)
			}
		}() //HACK
	} else {
		go func() {
			prevKeyStatus := hardwareKeyboard.Scan()
			for {
				keyStatus := hardwareKeyboard.Scan()
				time.Sleep(KEYPOLLINTERVALMS * time.Millisecond)
				if !keyStatus.Equal(prevKeyStatus) {
					prevKeyStatus = keyStatus
					if len(keyStatus) == 0 {
						CmdCh <- pulsegenui.CMDBTN_RELEASE //Changed and nothing pressed. report as released
					} else {
						for _, cmd := range keyStatus {
							CmdCh <- cmd
						}
					}
				}
			}
		}()
	}

	var oled BW128x64Display
	var errOled error

	if *pSimHw || *pSimDisplay {
		oled = &Fakedisplay128x64{}
	} else {
		//Display init
		fDisplay, errI2CHardware := os.OpenFile(I2CDEVICEFILE, os.O_RDWR, 0600)
		if errI2CHardware != nil {
			fmt.Printf("I2C Hardware open error %v", errI2CHardware.Error())
			os.Exit(-1)
		}

		var oledhw SSD1306_i2c
		oledhw, errOled = InitSSD1306_i2c(fDisplay, I2CADDRESS_OLED)
		if errOled != nil {
			fmt.Printf("I2C display init error %v", errOled.Error())
			os.Exit(-1)
		}
		oled = BW128x64Display(&oledhw)
	}

	go func() {
		for {
			oled.FullDisplayUpdate(BitmapToSSD1306Buffer(<-BitmapCh, false))
		}
	}()

	go func() {
		for {
			errSet := setPwm0Generator(<-ui.RfCmd)
			if errSet != nil {
				fmt.Printf("ERROR %v", errSet)
				os.Exit(-1)
			}
		}
	}()

	ui.Run() //This is where "business logic is ticking".
}
