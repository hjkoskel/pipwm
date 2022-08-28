module github.com/hjkoskel/pipwm/pulseGadget

go 1.18

require (
	github.com/hjkoskel/gomonochromebitmap v0.1.0-beta.1
	github.com/hjkoskel/govattu v0.1.0-beta.1
	github.com/hjkoskel/pipwm/pulsegenui v0.0.0-00010101000000-000000000000
	github.com/nsf/termbox-go v1.1.1
	github.com/stretchr/testify v1.8.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.3.4 // indirect
	golang.org/x/sys v0.0.0-20220825204002-c680a09ffe64 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hjkoskel/pipwm/pulsegenui => ../pulsegenui
