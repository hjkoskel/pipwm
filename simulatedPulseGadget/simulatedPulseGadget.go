/*
Simulated pulseGenerator gadget
Can be used for testing OLED ui rendering and PWM calcuations
*/
package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/hjkoskel/gomonochromebitmap"
	"github.com/hjkoskel/pipwm/pulseGenUi"

	"github.com/hjkoskel/gomonochromebitmap/gadgetSimUi"
)

const (
	BITMAP_BUTTONDIM = 126
	BITMAP_ROW0      = 300
	BITMAP_ROW1      = 443
	BITMAP_ROW2      = 587
	BITMAP_ROW3      = 730

	BITMAP_COL0 = 14
	BITMAP_COL1 = 162
	BITMAP_COL2 = 313
	BITMAP_COL3 = 462
)

const (
	BITMAP_BACKGROUNDFILE = "bgFuncGen.png"
)

const (
	BITMAP_DRAWBUTTONAREAS = false
)

func main() {
	sim := true

	BitmapCh := make(chan gomonochromebitmap.MonoBitmap, 1)
	CmdCh := make(chan string, 3) //Key pressesses. This gadget takes only one press per time
	ui := pulseGenUi.PulsGenUi{Bitmap: BitmapCh, Cmd: CmdCh, Simulate: sim}

	go ui.Run() //This is where "business logic is ticking". Leaving main loop for SDL library, it likes that

	//TODO load from file later
	reader, err := os.Open(BITMAP_BACKGROUNDFILE)
	if err != nil {
		fmt.Printf("Bitmap %v load failed err=%v\n", BITMAP_BACKGROUNDFILE, err.Error())
	}

	bgImage, _, err := image.Decode(reader)
	if err != nil {
		fmt.Printf("Bitmap %v decode failed err=%v\n", BITMAP_BACKGROUNDFILE, err.Error())
	}
	reader.Close() //reading is done?

	btnDimOnBitmap := gadgetSimUi.XyIntPair{X: BITMAP_BUTTONDIM, Y: BITMAP_BUTTONDIM} //Size on bitmap

	buttonDebugColor := color.RGBA{A: 255, R: 0, G: 255, B: 0}

	//Ok... I was lazy ass just copypasted and find&replaced this huge JSON here..
	gw := gadgetSimUi.GadgetWindow{
		MonoDisplays: []gadgetSimUi.MonochromeDisplay{
			gadgetSimUi.MonochromeDisplay{
				ID:          "",
				Corner:      gadgetSimUi.XyIntPair{X: 47, Y: 29},
				PixelSize:   gadgetSimUi.XyIntPair{X: 3, Y: 3},
				PixelGap:    gadgetSimUi.XyIntPair{X: 1, Y: 1},
				OnColor:     color.RGBA{A: 255, R: 200, G: 200, B: 0},
				OnColorDown: color.RGBA{A: 255, R: 0, G: 200, B: 200},
				OffColor:    color.RGBA{A: 255, R: 50, G: 50, B: 50},
				UpperRows:   8,
			},
		},
		Buttons: []gadgetSimUi.ButtonSettings{
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_7,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL0, Y: BITMAP_ROW0},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_8,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL1, Y: BITMAP_ROW0},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_9,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL2, Y: BITMAP_ROW0},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_UP,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL3, Y: BITMAP_ROW0},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},

			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_4,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL0, Y: BITMAP_ROW1},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_5,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL1, Y: BITMAP_ROW1},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_6,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL2, Y: BITMAP_ROW1},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_DOWN,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL3, Y: BITMAP_ROW1},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},

			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_1,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL0, Y: BITMAP_ROW2},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_2,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL1, Y: BITMAP_ROW2},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_3,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL2, Y: BITMAP_ROW2},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_BACK,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL3, Y: BITMAP_ROW2},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},

			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_ONOFF,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL0, Y: BITMAP_ROW3},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_0,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL1, Y: BITMAP_ROW3},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_DECIMALPOINT,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL2, Y: BITMAP_ROW3},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
			gadgetSimUi.ButtonSettings{
				ID:         pulseGenUi.CMDBTN_OK,
				Corner:     gadgetSimUi.XyIntPair{X: BITMAP_COL3, Y: BITMAP_ROW3},
				Dimensions: btnDimOnBitmap,
				DebugColor: buttonDebugColor,
				DebugEdges: BITMAP_DRAWBUTTONAREAS,
			},
		},
	}

	//SDL likes to be runned on main routine, not in goroutine
	gw.Initialize(bgImage)

	go func() {
		for {
			gw.ToDisplay <- gadgetSimUi.DisplayUpdate{Bitmap: <-BitmapCh, ID: ""} //Channel adaptor
		}
	}()

	go func() {
		for {
			arr := (<-gw.FromKeys).KeysDown
			if len(arr) == 0 {
				CmdCh <- ""
			} else {
				CmdCh <- arr[0] //Single press only
			}
		}
	}()

	//SDL wants to run as main goroutine
	runErr := gw.Run()
	if runErr != nil {
		fmt.Printf("Runtime error %v\n", runErr.Error())
	}
	gw.Quit()

}
