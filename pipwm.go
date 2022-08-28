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
	"os"
	"time"

	"github.com/hjkoskel/govattu"
)

var hw govattu.Vattu

func setPwm0Generator(hi float64, lo float64, verbose bool) error {
	if !(0 < hi) { //SHUTDOWN TO LO
		hw.PwmSetMode(false, false, false, false)
		hw.PinMode(18, govattu.ALToutput)
		hw.PinClear(18)
		return nil
	}

	if !(0 < lo) { //SHUTDOWN TO HI
		hw.PwmSetMode(false, false, false, false)
		hw.PinMode(18, govattu.ALToutput)
		hw.PinSet(18)
		return nil
	}

	pulseTiming := govattu.RfPulseSettings{On: time.Duration(time.Nanosecond * time.Duration(hi*1000.0)), Off: time.Duration(time.Nanosecond * time.Duration(lo*1000.0))}
	pulseSettings, setErr := pulseTiming.GetSettings()
	if setErr != nil {
		return setErr
	}

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

	hw.PinMode(18, govattu.ALT5)            //ALT5 function for 18 is PWM0
	hw.PwmSetMode(true, true, false, false) // PwmSetMode(en0 bool, ms0 bool, en1 bool, ms1 bool)   enable and set to mark-space mode
	hw.PwmSetClock(uint32(pulseSettings.Pwmc))
	hw.Pwm0SetRange(uint32(pulseSettings.Pwmr))
	hw.Pwm0Set(uint32(pulseSettings.Pwm))
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

	if pSimulated {
		hw = &govattu.DoNothingPi{}
	} else {
		var errhw error
		hw, errhw = govattu.Open()
		if errhw != nil {
			fmt.Printf("\n\n%v\n", errhw.Error())
			os.Exit(-1)
		}
		defer hw.Close()
	}

	if (0 <= pLo) && (0 <= pHi) {
		err := setPwm0Generator(pHi, pLo, pVerbose)
		if err != nil {
			fmt.Printf("ERR %s\n", err.Error())
			os.Exit(-1)
		}
		return
	}

	if (0 <= pPercent) && (pPercent <= 100) { //Percent is defined try calculating missing info
		ratio := pPercent / 100.0 //Simpler
		if 0 <= pTr {             //period and percent defined
			err := setPwm0Generator(pTr*ratio, pTr*(1.0-ratio), pVerbose)
			if err != nil {
				fmt.Printf("ERR=%v\n", err.Error())
				os.Exit(-1)
			}
			return
		}
		if (pLo < 0) && (!(pHi < 0)) { //Hi time is known, but not lo... so calc period
			if pVerbose {
				fmt.Printf("Calculate by using hi time and percent\n")
			}
			err := setPwm0Generator(pHi, pHi*(1/ratio-1), pVerbose)
			if err != nil {
				fmt.Printf("ERR=%v\n", err.Error())
				os.Exit(-1)
			}
			return
		}
		if (pHi < 0) && (!(pLo < 0)) { //Lo time is known, but not hi... so calc period
			if pVerbose {
				fmt.Printf("Calculate by using lo time and percent\n")
			}
			err := setPwm0Generator(ratio*pLo/(1-ratio), pLo, pVerbose)
			if err != nil {
				fmt.Printf("ERR=%v\n", err.Error())
				os.Exit(-1)
			}
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
			err := setPwm0Generator(pHi, pTr-pHi, pVerbose)
			if err != nil {
				fmt.Printf("ERR=%v\n", err.Error())
				os.Exit(-1)
			}
			return
		}
		if 0 <= pLo {
			if pTr < pLo {
				fmt.Printf("ERROR period repeat time %v < low time %v\n", pTr, pLo)
				os.Exit(-1)
			}
			err := setPwm0Generator(pTr-pLo, pLo, pVerbose)
			if err != nil {
				fmt.Printf("ERR=%v\n", err.Error())
				os.Exit(-1)
			}
			return
		}
	}

	if -1000 < pSa {
		//USE SERVO MODE
		pHi = (pSa/90)*pSs + pSc
		if pTr <= 0 { //Use default period
			pTr = 20000 //20ms period for servos
		}
		err := setPwm0Generator(pHi, pTr-pHi, pVerbose)
		if err != nil {
			fmt.Printf("ERR=%v\n", err.Error())
			os.Exit(-1)
		}
		return
	}
	fmt.Printf("Not enough settings or want shut down PWM\n")
	setPwm0Generator(0, 0, pVerbose) //Leave low
}
