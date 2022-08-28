/*
Embedded function generator with display and keys
Produces 128x64 bitmap and recieves commands from 4x4 keypad. Virtual or real

*/
package pulsegenui

import (
	"fmt"
	"image"
	"math"
	"strconv"
	"time"

	"github.com/hjkoskel/gomonochromebitmap"
	"github.com/hjkoskel/govattu"
)

//For not using magic strings
const (
	CMDBTN_0            = "0"
	CMDBTN_1            = "1"
	CMDBTN_2            = "2"
	CMDBTN_3            = "3"
	CMDBTN_4            = "4"
	CMDBTN_5            = "5"
	CMDBTN_6            = "6"
	CMDBTN_7            = "7"
	CMDBTN_8            = "8"
	CMDBTN_9            = "9"
	CMDBTN_ONOFF        = "onoff"
	CMDBTN_DECIMALPOINT = "."
	CMDBTN_UP           = "up"
	CMDBTN_DOWN         = "down"
	CMDBTN_BACK         = "back"
	CMDBTN_OK           = "ok"
	CMDBTN_RELEASE      = ""
)

type UiPageNumber int

const (
	PAGE_REGS = iota
	PAGE_LOHI
	PAGE_SERVO
	PAGE_MAXCOUNT //THE LAST
)

//Page main data separated. Each page know hows how to render

type RegsListItem int

const (
	LISTITEM_REG_PWMC = iota
	LISTITEM_REG_PWMR
	LISTITEM_REG_PWM
	LISTITEM_REG_MAXCOUNT //THE LAST
)

type UiPageRegs struct {
	Registers govattu.RfSettings

	SelectedField RegsListItem
	Edit          bool //active editing
	EditValue     uint32
}

func (p *UiPageRegs) StartEdit() {
	p.Edit = true
	switch p.SelectedField {
	case LISTITEM_REG_PWMC:
		p.EditValue = p.Registers.Pwmc
	case LISTITEM_REG_PWMR:
		p.EditValue = p.Registers.Pwmr
	case LISTITEM_REG_PWM:
		p.EditValue = p.Registers.Pwm
	}
}
func (p *UiPageRegs) StoreEdit() {
	p.Edit = false
	switch p.SelectedField {
	case LISTITEM_REG_PWMC:
		p.Registers.Pwmc = p.EditValue
	case LISTITEM_REG_PWMR:
		p.Registers.Pwmr = p.EditValue
	case LISTITEM_REG_PWM:
		p.Registers.Pwm = p.EditValue
	}
}

const (
	LISTITEM_LOHI_SCALE_NANO = iota
	LISTITEM_LOHI_SCALE_MICRO
	LISTITEM_LOHI_SCALE_MILLI
	LISTITEM_LOHI_SCALE_UNIT
	LISTITEM_LOHI_HI
	LISTITEM_LOHI_LO
	LISTITEM_LOHI_MAXCOUNT
)

const (
	SCALETIME_NANO  = 1
	SCALETIME_MICRO = 1000
	SCALETIME_MILLI = 1000000
	SCALETIME_UNIT  = 1000000000
)

type UiPageLoHi struct {
	Scale         int64 // 1=nano, 1000=micro 1000 000=milli,  1000 000 000 s
	Lo            float64
	Hi            float64
	SelectedField int
	Edit          bool
	EditedValue   string //Decimal points need to store
}

func (p *UiPageLoHi) StartEdit() {
	switch p.SelectedField { //Only allowed fields
	case LISTITEM_LOHI_LO:
		p.Edit = true
		p.EditedValue = "" //Lets start with empty
	case LISTITEM_LOHI_HI:
		p.Edit = true
		p.EditedValue = ""
	}
}

func (p *UiPageLoHi) StoreEdit() {
	p.Edit = false
	newValue, err := strconv.ParseFloat(p.EditedValue, 32)
	if err != nil {
		return //invalid
	}

	switch p.SelectedField { //Only allowed fields
	case LISTITEM_LOHI_LO:
		p.Lo = newValue
	case LISTITEM_LOHI_HI:
		p.Hi = newValue
	}
}

type UiPageServo struct {
	Angle float32
}

func (p *UiPageServo) LimitAngle(minAngle float32, maxAngle float32) {
	p.Angle = float32(math.Max(float64(minAngle), math.Min(float64(maxAngle), float64(p.Angle))))
}

type PulseHardwareCommand struct {
	OutputEnabled bool
	Rf            govattu.RfPulseSettings
}

func (a PulseHardwareCommand) String() string {
	if a.OutputEnabled {
		return fmt.Sprintf("ENABLED on:%s off:%s", a.Rf.On, a.Rf.Off)
	}
	return fmt.Sprintf("DISABLE (on:%s off:%s)", a.Rf.On, a.Rf.Off)
}

type PulsGenUi struct {
	//Simulate bool //Important when testing
	Bitmap chan gomonochromebitmap.MonoBitmap
	Cmd    chan string //Key pressesses. This gadget takes only one press per time
	RfCmd  chan PulseHardwareCommand

	//---- status composed to status bar -----
	SignalOut bool //Is feeding actively
	//--- Sub page status ---
	Page UiPageNumber

	StatusRegs  UiPageRegs
	StatusLoHi  UiPageLoHi
	StatusServo UiPageServo

	//--- Selected fonts
	titleFont  map[rune]gomonochromebitmap.MonoBitmap
	normalFont map[rune]gomonochromebitmap.MonoBitmap
}

func (p *PulsGenUi) toggleOutputOnOff() {
	p.SignalOut = !p.SignalOut
	p.setHardwareFromRegisters()
}

func (p *PulsGenUi) readKey() string {
	cmd := <-p.Cmd
	if cmd == CMDBTN_ONOFF {
		p.toggleOutputOnOff()
	}

	return cmd
}

//This is main place where values are set

func (p *PulsGenUi) setHardwareFromRegisters() {
	//if !p.Simulate {
	if p.StatusRegs.Registers.IsOffline() {
		p.SignalOut = false //Disable output
	}

	p.RfCmd <- PulseHardwareCommand{OutputEnabled: p.SignalOut, Rf: p.StatusRegs.Registers.GetTiming()}
	//}
}

func (p *PulsGenUi) setHardwareFromLoHi() error {
	//fmt.Printf("TODO SET HARDWARE FROM LO:%f ns HI:%f ns\n",p.StatusLoHi.Lo,p.StatusLoHi.Hi)
	d := p.StatusLoHi
	pt := govattu.RfPulseSettings{
		On:  time.Duration(time.Nanosecond * time.Duration(d.Hi*float64(d.Scale))),
		Off: time.Duration(time.Nanosecond * time.Duration(d.Lo*float64(d.Scale))),
	}
	var err error
	p.StatusRegs.Registers, err = pt.GetSettings()
	if err != nil {
		return err
	}
	p.setHardwareFromRegisters()
	return err
}

//Idea of chained settings
func (p *PulsGenUi) setHardwareFromServo() error {
	p.StatusLoHi.Hi = float64((p.StatusServo.Angle/90)*0.5 + 1.5)
	p.StatusLoHi.Lo = (20 - p.StatusLoHi.Hi) //20ms
	p.StatusLoHi.Scale = SCALETIME_MILLI
	return p.setHardwareFromLoHi()
}

func (p *PulsGenUi) initializeBitmapWithHeader() gomonochromebitmap.MonoBitmap {
	result := gomonochromebitmap.NewMonoBitmap(128, 64, false)
	a := image.Rectangle{}
	a.Min.X = 0
	a.Max.X = 126
	a.Min.Y = 0
	a.Max.Y = 8

	// sec  millisec microsec nanosec
	pulseTiming := p.StatusRegs.Registers.GetTiming()
	up := float32(pulseTiming.On.Nanoseconds())
	down := float32(pulseTiming.Off.Nanoseconds())

	//fmt.Printf("Output pulse is %v ns up and %v ns down\n", up, down)

	headerText := ""
	if !p.StatusRegs.Registers.IsOffline() {
		headerText = fmt.Sprintf("%.3f/%.3fns", up, down)

		divid := float32(1000)
		if ((divid < up) || (up <= 0.00001)) && ((divid < down) || (down <= 0.00001)) {
			headerText = fmt.Sprintf("%.3f/%.3fus", up/divid, down/divid)
		}

		divid = float32(1000 * 1000)
		if ((divid < up) || (up <= 0.00001)) && ((divid < down) || (down <= 0.00001)) {
			headerText = fmt.Sprintf("%.3f/%.3fms", up/divid, down/divid)
		}

		divid = float32(1000 * 1000 * 1000)
		if ((divid < up) || (up <= 0.00001)) && ((divid < down) || (down <= 0.00001)) {
			headerText = fmt.Sprintf("%.3f/%.3fs", up/divid, down/divid)
		}
	}

	if p.SignalOut {
		headerText = "ON " + headerText
	} else {
		headerText = "OFF " + headerText
	}

	result.Print(headerText, p.titleFont, 1, 0, a, true, true, false, false)
	return result
}

/*
Keeps running
Is actively updating only by input

TODO capture
*/
func (p *PulsGenUi) Run() {
	p.titleFont = gomonochromebitmap.GetFont_5x7()
	p.normalFont = gomonochromebitmap.GetFont_8x8()

	p.Page = PAGE_REGS
	for {
		p.render()
		//Move in menu
		switch p.readKey() {
		case CMDBTN_UP:
			if 0 < p.Page {
				p.Page = (p.Page - 1) % PAGE_MAXCOUNT
			}
		case CMDBTN_DOWN:
			p.Page = (p.Page + 1) % PAGE_MAXCOUNT
		case CMDBTN_OK: //Picks item
			switch p.Page {
			case PAGE_REGS:
				p.runRegs()
			case PAGE_LOHI:
				p.runLoHi()
			case PAGE_SERVO:
				p.runServo()
			}
		}
	}
}

func (p *PulsGenUi) render() {
	bm := p.initializeBitmapWithHeader()

	a := image.Rectangle{}
	a.Min.X = 4
	a.Max.X = 126
	a.Max.Y = 64

	a.Min.Y = 15
	bm.Print("Regs", p.normalFont, 1, 0, a, true, true, p.Page == PAGE_REGS, false)
	a.Min.Y += 10
	bm.Print("Lo-Hi", p.normalFont, 1, 0, a, true, true, p.Page == PAGE_LOHI, false)
	a.Min.Y += 10
	bm.Print("Servo", p.normalFont, 1, 0, a, true, true, p.Page == PAGE_SERVO, false)
	p.Bitmap <- bm
}

/*
returns value, done and
When there are no decimals, it is possible to work with int, not string
*/
func editUint32Iteration(value uint32, cmd string, minValue uint32, maxValue uint32) uint32 {
	switch cmd {
	case CMDBTN_UP:
		if value < maxValue {
			return value + 1
		}
	case CMDBTN_DOWN:
		if minValue < value {
			return value - 1
		}
	default:
		//Parse if number
		i, convErr := strconv.ParseInt(cmd, 10, 32)
		if convErr == nil {
			v := float64(value*10 + uint32(i))
			return uint32(math.Min(math.Max(float64(minValue), v), float64(maxValue)))
		}
	}
	return value
}

const FLOATFORMATSTRING = "%.5f"

func editFloatIteration(value string, cmd string, maxChars int) string {
	if maxChars <= len(value) {
		return value
	}
	switch cmd { //Yes, this is so dumb :D
	case CMDBTN_0, CMDBTN_1, CMDBTN_2, CMDBTN_3, CMDBTN_4, CMDBTN_5, CMDBTN_6, CMDBTN_7, CMDBTN_8, CMDBTN_9, CMDBTN_DECIMALPOINT:
		return value + cmd
	}
	return value
}

func (p *PulsGenUi) runRegs() {
	p.StatusRegs.Edit = false
	p.StatusRegs.EditValue = 0
	cmdGiven := ""
	for cmdGiven != CMDBTN_BACK { //TODO NO MAGIC STRINGS
		p.renderRegs()
		cmdGiven = p.readKey()
		switch cmdGiven {
		case CMDBTN_UP:
			if 0 < p.StatusRegs.SelectedField {
				p.StatusRegs.SelectedField = (p.StatusRegs.SelectedField - 1) % LISTITEM_REG_MAXCOUNT
			}
		case CMDBTN_DOWN:
			p.StatusRegs.SelectedField = (p.StatusRegs.SelectedField + 1) % LISTITEM_REG_MAXCOUNT
		case CMDBTN_OK: //Going to edit
			cmdGiven = ""
			p.StatusRegs.StartEdit()
			p.StatusRegs.EditValue = 0
			for (cmdGiven != CMDBTN_OK) && (cmdGiven != CMDBTN_BACK) {
				p.renderRegs()
				cmdGiven = p.readKey()
				switch p.StatusRegs.SelectedField {
				case LISTITEM_REG_PWM:
					p.StatusRegs.EditValue = editUint32Iteration(p.StatusRegs.EditValue, cmdGiven, 0, 0xFFFF)
				case LISTITEM_REG_PWMC:
					p.StatusRegs.EditValue = editUint32Iteration(p.StatusRegs.EditValue, cmdGiven, govattu.MINPWMC, govattu.MAXPWMC)
				case LISTITEM_REG_PWMR:
					p.StatusRegs.EditValue = editUint32Iteration(p.StatusRegs.EditValue, cmdGiven, 0, 0xFFFF)
				}
			}
			if cmdGiven == CMDBTN_OK {
				//fmt.Printf("Storing edited value %v", p.StatusRegs.EditValue)
				p.StatusRegs.StoreEdit()
				p.setHardwareFromRegisters()
			}
			p.StatusRegs.Edit = false
		}
	}
}

func (p *PulsGenUi) renderRegs() {
	a := image.Rectangle{}
	bm := p.initializeBitmapWithHeader()
	a.Min.X = 0
	a.Max.X = 126
	a.Min.Y = 16
	a.Max.Y = 16 + 8
	d := p.StatusRegs
	if (d.Edit) && (d.SelectedField == LISTITEM_REG_PWMC) {
		bm.Print(fmt.Sprintf("PWMC: %v_", d.EditValue), p.normalFont, 1, 1, a, true, true, true, false)
	} else {
		bm.Print(fmt.Sprintf("PWMC: %v", d.Registers.Pwmc), p.normalFont, 1, 1, a, true, true, (d.SelectedField == LISTITEM_REG_PWMC), false)
	}

	a.Min.Y = 2 * 16
	a.Max.Y = 2*16 + 8
	if (d.Edit) && (d.SelectedField == LISTITEM_REG_PWMR) {
		bm.Print(fmt.Sprintf("PWMR: %v_", d.EditValue), p.normalFont, 1, 1, a, true, true, true, false)
	} else {
		bm.Print(fmt.Sprintf("PWMR: %v", d.Registers.Pwmr), p.normalFont, 1, 1, a, true, true, (d.SelectedField == LISTITEM_REG_PWMR), false)
	}

	a.Min.Y = 3 * 16
	a.Max.Y = 3*16 + 8
	if (d.Edit) && (d.SelectedField == LISTITEM_REG_PWM) {
		bm.Print(fmt.Sprintf("PWM: %v_", d.EditValue), p.normalFont, 1, 1, a, true, true, true, false)
	} else {
		bm.Print(fmt.Sprintf("PWM: %v", d.Registers.Pwm), p.normalFont, 1, 1, a, true, true, (d.SelectedField == LISTITEM_REG_PWM), false)
	}
	p.Bitmap <- bm
}

func (p *PulsGenUi) runLoHi() {
	cmdGiven := ""
	if p.StatusLoHi.Scale == 0 {
		p.StatusLoHi.Scale = SCALETIME_MICRO
	}
	for cmdGiven != CMDBTN_BACK {
		p.renderLoHi()
		cmdGiven = p.readKey()
		switch cmdGiven {
		case CMDBTN_UP:
			if 0 < p.StatusLoHi.SelectedField {
				p.StatusLoHi.SelectedField = (p.StatusLoHi.SelectedField - 1) % LISTITEM_LOHI_MAXCOUNT
			}
		case CMDBTN_DOWN:
			p.StatusLoHi.SelectedField = (p.StatusLoHi.SelectedField + 1) % LISTITEM_LOHI_MAXCOUNT
		case CMDBTN_OK: //Going to edit
			switch p.StatusLoHi.SelectedField {
			case LISTITEM_LOHI_SCALE_NANO:
				p.StatusLoHi.Scale = SCALETIME_NANO
				p.setHardwareFromLoHi()
			case LISTITEM_LOHI_SCALE_MICRO:
				p.StatusLoHi.Scale = SCALETIME_MICRO
				p.setHardwareFromLoHi()
			case LISTITEM_LOHI_SCALE_MILLI:
				p.StatusLoHi.Scale = SCALETIME_MILLI
				p.setHardwareFromLoHi()
			case LISTITEM_LOHI_SCALE_UNIT:
				p.StatusLoHi.Scale = SCALETIME_UNIT
				p.setHardwareFromLoHi()

			case LISTITEM_LOHI_HI, LISTITEM_LOHI_LO:
				p.StatusLoHi.StartEdit()
				cmdGiven = ""
				for (cmdGiven != CMDBTN_OK) && (cmdGiven != CMDBTN_BACK) {
					p.renderLoHi()
					cmdGiven = p.readKey()
					p.StatusLoHi.EditedValue = editFloatIteration(p.StatusLoHi.EditedValue, cmdGiven, 10)
				}
				if cmdGiven == CMDBTN_OK {
					p.StatusLoHi.StoreEdit()
					p.setHardwareFromLoHi()
				}
				p.StatusRegs.Edit = false //Like cancel with back
			}
		}
	}
}
func (p *PulsGenUi) renderLoHi() {
	a := image.Rectangle{}
	bm := p.initializeBitmapWithHeader()
	a.Min.X = 5
	a.Min.Y = 9
	a.Max.Y = 64
	a.Max.X = 126
	d := p.StatusLoHi

	bm.Print("n", p.normalFont, 0, 1, a, true, true, (d.Scale == SCALETIME_NANO) || (d.SelectedField == LISTITEM_LOHI_SCALE_NANO), false)
	a.Min.X += 20
	bm.Print("u", p.normalFont, 0, 1, a, true, true, (d.Scale == SCALETIME_MICRO) || (d.SelectedField == LISTITEM_LOHI_SCALE_MICRO), false)
	a.Min.X += 20
	bm.Print("m", p.normalFont, 0, 1, a, true, true, (d.Scale == SCALETIME_MILLI) || (d.SelectedField == LISTITEM_LOHI_SCALE_MILLI), false)
	a.Min.X += 20
	bm.Print("sec", p.normalFont, 0, 1, a, true, true, (d.Scale == SCALETIME_UNIT) || (d.SelectedField == LISTITEM_LOHI_SCALE_UNIT), false)

	unitname := ""
	switch d.Scale {
	case SCALETIME_NANO:
		unitname = "nanoseconds"
	case SCALETIME_MICRO:
		unitname = "microseconds"
	case SCALETIME_MILLI:
		unitname = "milliseconds"
	case SCALETIME_UNIT:
		unitname = "seconds"
	}
	a.Min.X = 2
	a.Min.Y = 55
	bm.Print(unitname, p.normalFont, 0, 1, a, true, true, false, false)

	a.Min.X = 3
	a.Min.Y = 22
	if d.Edit && (d.SelectedField == LISTITEM_LOHI_HI) {
		bm.Print(d.EditedValue+"_", p.normalFont, 0, 0, a, true, true, d.SelectedField == LISTITEM_LOHI_HI, false)
	} else {
		bm.Print(fmt.Sprintf(FLOATFORMATSTRING, d.Hi), p.normalFont, 0, 0, a, true, true, d.SelectedField == LISTITEM_LOHI_HI, false)
	}

	//Percent display
	if 0 < (d.Lo + d.Hi) {
		ratio := (d.Hi) / (d.Lo + d.Hi)
		ratioPixel := int(math.Floor(127.0 * ratio))
		bm.Hline(0, ratioPixel, 33, true)
		bm.Vline(ratioPixel, 33, 40, true)
		bm.Hline(ratioPixel, 127, 40, true)

		a.Min.X = 3
		a.Min.Y = 48
		//bm.Print(text, font, lineSpacing, gap, area, drawTrue, drawFalse, invert, wrap)
		s := fmt.Sprintf("%.2f%%", 100*ratio)
		if 0.99 < ratio {
			s = "100%"
		}
		bm.Print(s, p.normalFont, 0, 0, a, true, true, false, false)
	}

	a.Min.X = 60
	a.Min.Y = 48
	if d.Edit && (d.SelectedField == LISTITEM_LOHI_LO) {
		bm.Print(d.EditedValue+"_", p.normalFont, 0, 0, a, true, true, d.SelectedField == LISTITEM_LOHI_LO, false)
	} else {
		bm.Print(fmt.Sprintf(FLOATFORMATSTRING, d.Lo), p.normalFont, 0, 0, a, true, true, d.SelectedField == LISTITEM_LOHI_LO, false)
	}

	p.Bitmap <- bm
}

/*
Servo have pre-set settings
=0 center
1=-90
2=-70
3=-50

9=90
*/

const (
	SERVOPRESET_0 = 0.0
	SERVOPRESET_1 = -90
	SERVOPRESET_2 = -80
	SERVOPRESET_3 = -50
	SERVOPRESET_4 = -20
	SERVOPRESET_5 = 20
	SERVOPRESET_6 = 50
	SERVOPRESET_7 = 80
	SERVOPRESET_8 = 85
	SERVOPRESET_9 = -90
)
const (
	SERVOSWINGTIME = 3500 //How many milliseconds to edge to edge
	SERVOMINANGLE  = -90.0
	SERVOMAXANGLE  = 90.0
)

func (p *PulsGenUi) getServoAngleFromRelease() float32 {
	cmdGiven := "aaa"
	tDown := time.Now()
	for cmdGiven != CMDBTN_RELEASE {
		cmdGiven = p.readKey()
	}
	tUp := time.Now()
	return 180 * float32(float32(tUp.Sub(tDown).Nanoseconds())/float32(1000*1000)) / float32(SERVOSWINGTIME)
}

func (p *PulsGenUi) runServo() {
	p.StatusServo.Angle = 0
	cmdGiven := ""
	for cmdGiven != CMDBTN_BACK { //TODO NO MAGIC STRINGS
		switch cmdGiven {
		case CMDBTN_UP:
			p.StatusServo.Angle += p.getServoAngleFromRelease()
		case CMDBTN_DOWN:
			p.StatusServo.Angle -= p.getServoAngleFromRelease()
		}
		p.StatusServo.LimitAngle(SERVOMINANGLE, SERVOMAXANGLE)
		p.renderServo()
		p.setHardwareFromServo()
		cmdGiven = p.readKey()
	}
}

const (
	GAUGEY      = 53
	GAUGEHEIGHT = 35 + 5
)

func (p *PulsGenUi) renderServo() {
	bm := p.initializeBitmapWithHeader()

	bm.Circle(image.Point{X: 64, Y: GAUGEY}, GAUGEHEIGHT, true)

	a := image.Rectangle{}
	a.Min.X = 0
	a.Min.Y = GAUGEY
	a.Max.X = 127
	a.Max.Y = 64

	bm.Fill(a, false)
	a.Min.X = 5
	a.Min.Y += 2
	bm.Print(fmt.Sprintf("%.2f degree", p.StatusServo.Angle), p.normalFont, 0, 0, a, true, false, false, false)

	rad := math.Pi * float64(p.StatusServo.Angle-90) / (180.0)

	bm.Line(image.Point{X: 64, Y: GAUGEY}, image.Point{
		X: 64 + int(math.Cos(rad)*float64(GAUGEHEIGHT)),
		Y: GAUGEY + int(math.Sin(rad)*float64(GAUGEHEIGHT)),
	}, true)

	p.Bitmap <- bm
}
