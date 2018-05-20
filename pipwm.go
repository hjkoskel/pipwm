/*
Program for controlling raspberry pwm
Sets HARDWARE pwm to new setpoint and makes exit.

At first only one pin is supported. Going to extend if audience requires. One pin is enough for my own use.

I use this for testing my to "be ordered" high voltage high frequency rf driver

Acts also as example project for govattu

Later adding PWM1 and software PWM is also option
*/

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/hjkoskel/govattu"
)

const PWMCLOCKNS float64 = 52.0833333333 //=(1000*1000*1000)/19200000  104.1666667 / 2.0 //One "block" takes this long
const MINPWMC int16 = 2
const MAXPWMC int16 = 2090

type RfSettings struct {
	Pwmc int16
	Pwmr int16
	Pwm  int16
}

func (p *RfSettings) Equal(a RfSettings) bool {
	return (p.Pwm == a.Pwm) && (p.Pwmr == a.Pwmr) && (p.Pwmc == p.Pwmc)
}

type RfPulseSettings struct {
	On  time.Duration
	Off time.Duration
}

func (p *RfPulseSettings) Equal(a RfPulseSettings) bool {
	return (p.On == a.On) && (p.Off == a.Off)
}

//Naive solution
func (p *RfPulseSettings) GetSettings() RfSettings {
	onns := p.On.Nanoseconds()                                    //On time in microseconds
	periodns := float64(p.On.Nanoseconds() + p.Off.Nanoseconds()) //Period lenght in microseconds
	result := RfSettings{}
	clocksPerPeriod := periodns / PWMCLOCKNS //How many times PWM clock will clock on period
	//Pick minimal divider = largest possible frequency
	result.Pwmc = int16(math.Ceil(float64(clocksPerPeriod)/float64(MAXPWMC))) + 1
	if result.Pwmc < MINPWMC {
		result.Pwmc = MINPWMC
	}
	result.Pwmr = int16(float64(clocksPerPeriod) / float64(result.Pwmc))
	//Now calculate how many points are required for period
	blockTime := int64(PWMCLOCKNS * float64(result.Pwmc)) //WRONG   Now this tells how long one block takes
	result.Pwm = int16(onns / blockTime)
	return result
}

func (p *RfSettings) GetTiming() RfPulseSettings {
	result := RfPulseSettings{}
	blockTime := PWMCLOCKNS * float64(p.Pwmc)
	result.On = time.Duration(blockTime*float64(p.Pwm)) * time.Nanosecond
	result.Off = time.Duration(blockTime*float64(p.Pwmr-p.Pwm)) * time.Nanosecond
	return result
}

func setPwm0Generator(simulated bool, hi float64, lo float64, verbose bool) error {
	if !(0 < hi) { //SHUTDOWN TO LO
		if simulated {
			if verbose {
				fmt.Printf("SIMULATED: Shutdown to low")
			}
		} else {
			govattu.PwmSetMode(false, false, false, false)
			govattu.PinMode(18, govattu.ALToutput)
			govattu.PinClear(18)
		}
		return nil
	}

	if !(0 < lo) { //SHUTDOWN TO HI
		if simulated {
			if verbose {
				fmt.Printf("SIMULATED: Shutdown to hi")
			}
		} else {
			govattu.PwmSetMode(false, false, false, false)
			govattu.PinMode(18, govattu.ALToutput)
			govattu.PinSet(18)
		}
		return nil
	}

	pulseTiming := RfPulseSettings{On: time.Duration(time.Nanosecond * time.Duration(hi*1000.0)), Off: time.Duration(time.Nanosecond * time.Duration(lo*1000.0))}
	pulseSettings := pulseTiming.GetSettings()
	actualTiming := pulseSettings.GetTiming()

	if verbose {
		fmt.Printf("Pulse settings are %#v\n", pulseSettings)
		fmt.Printf("Requested on/off  %v/%v us  got %v/%v us\n",
			float64(pulseTiming.On.Nanoseconds())/1000.0, float64(pulseTiming.Off.Nanoseconds())/1000.0,
			float64(actualTiming.On.Nanoseconds())/1000.0, float64(actualTiming.Off.Nanoseconds())/1000.0)
	}
	if pulseSettings.Pwmr <= pulseSettings.Pwm {
		err := fmt.Errorf("Cant do requested on/off  %v/%v us pulse", float64(pulseTiming.On.Nanoseconds())/1000.0, float64(pulseTiming.Off.Nanoseconds())/1000.0)
		fmt.Printf("\n\n%v\n", err.Error())
		return err
	}

	if 4095 < pulseSettings.Pwmc {
		err := fmt.Errorf("Cant do requested on/off  %v/%v us pulse pwmc can not be divided more", float64(pulseTiming.On.Nanoseconds())/1000.0, float64(pulseTiming.Off.Nanoseconds())/1000.0)
		fmt.Printf("\n\n%v\n", err.Error())
		return err
	}

	//Lets access hardware
	if !simulated {
		govattu.PinMode(18, govattu.ALT5)            //ALT5 function for 18 is PWM0
		govattu.PwmSetMode(true, true, false, false) // PwmSetMode(en0 bool, ms0 bool, en1 bool, ms1 bool)   enable and set to mark-space mode
		govattu.PwmSetClock(uint32(pulseSettings.Pwmc))
		govattu.Pwm0SetRange(uint32(pulseSettings.Pwmr))
		govattu.Pwm0Set(uint32(pulseSettings.Pwm))
	}
	return nil
}

func main() {
	ppVerbose := flag.Bool("v", false, "verbose printout")
	ppHi := flag.Float64("hi", -1, "hi microseconds TODO range")
	ppLo := flag.Float64("lo", -1, "lo microseconds TODO range")
	ppTr := flag.Float64("tr", -1, "time repeat (period) alternative lo or hi")
	ppPercent := flag.Float64("p", -1, "pulse ration in percent 0-100")
	ppSa := flag.Float64("sa", -1000, "servo angle  -90 to 90 is nominal wider is accepted also")
	ppSc := flag.Float64("sc", 1500, "servo center in microseconds")
	ppSs := flag.Float64("ss", 500, "servo swing in microseconds")

	ppSimulated := flag.Bool("sim", false, "simulated non rasperry platform")
	//ppKeypad := flag.Bool("keypad", false, "Physical 4x4 keymatrix on GPIO")
	//TODO OLED OSOITE ppOled:= flag.Int("oled", "", usage)

	flag.Parse()

	pVerbose := *ppVerbose
	pHi := *ppHi
	pLo := *ppLo
	pTr := *ppTr
	pPercent := *ppPercent
	pSa := *ppSa
	pSc := *ppSc
	pSs := *ppSs
	pSimulated := *ppSimulated

	if !pSimulated {
		hw, err := govattu.Open()
		if err != nil {
			fmt.Printf("\n\n%v\n", err.Error())
			os.Exit(-1)
		}
		defer hw.Close()
	}

	if (0 <= pLo) && (0 <= pHi) {
		setPwm0Generator(pSimulated, pHi, pLo, pVerbose)
		return
	}

	if (0 <= pPercent) && (pPercent <= 100) { //Percent is defined try calculating missing info
		ratio := pPercent / 100.0 //Simpler
		if 0 <= pTr {             //period and percent defined
			setPwm0Generator(pSimulated, pTr*ratio, pTr*(1.0-ratio), pVerbose)
			return
		}
		if (pLo < 0) && (!(pHi < 0)) { //Hi time is known, but not lo... so calc period
			if pVerbose {
				fmt.Printf("Calculate by using hi time and percent\n")
			}
			setPwm0Generator(pSimulated, pHi, pHi*(1/ratio-1), pVerbose)
			return
		}
		if (pHi < 0) && (!(pLo < 0)) { //Lo time is known, but not hi... so calc period
			if pVerbose {
				fmt.Printf("Calculate by using lo time and percent\n")
			}
			setPwm0Generator(pSimulated, ratio*pLo/(1-ratio), pLo, pVerbose)
			return
		}
	}

	//Time repeat is defined
	if 0 < pTr {
		if 0 <= pHi {
			if pTr < pHi {
				fmt.Printf("ERROR period repeat time %v < high time %v\n", pTr, pHi)
				os.Exit(-1)
			}
			setPwm0Generator(pSimulated, pHi, pTr-pHi, pVerbose)
			return
		}
		if 0 <= pLo {
			if pTr < pLo {
				fmt.Printf("ERROR period repeat time %v < low time %v\n", pTr, pLo)
				os.Exit(-1)
			}
			setPwm0Generator(pSimulated, pTr-pLo, pLo, pVerbose)
			return
		}
	}

	if -1000 < pSa {
		//USE SERVO MODE
		pHi = (pSa/90)*pSs + pSc
		if pTr <= 0 { //Use default period
			pTr = 20000 //20ms period for servos
		}
		setPwm0Generator(pSimulated, pHi, pTr-pHi, pVerbose)
		return
	}
	fmt.Printf("Not enough settings or want shut down PWM\n")
	setPwm0Generator(pSimulated, 0, 0, pVerbose) //Leave low
}
