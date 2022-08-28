module simulatedPulseGadget

go 1.18

require (
	github.com/hjkoskel/gomonochromebitmap v0.1.0-beta.1
	github.com/hjkoskel/govattu v0.1.0-beta.1
	github.com/hjkoskel/pipwm/pulsegenui v0.0.0-00010101000000-000000000000
)

require (
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/veandco/go-sdl2 v0.4.25 // indirect
	golang.org/x/sys v0.0.0-20220825204002-c680a09ffe64 // indirect
)

replace github.com/hjkoskel/pipwm/pulsegenui => ../pulsegenui
