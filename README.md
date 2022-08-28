pipwm
=====

The pipwm is command line utility for controlling hardware PWM0 pin (18) on raspberry
At this first this is used only as test tool for my other hardware project
This is also good example how to use https://github.com/hjkoskel/govattu library on go

## Compiling ##

```sh
GOARCH=arm GOOS=linux go build
```
should cross compile source code to raspberry binary

## Usage ##
Due the limitations of govattu, this program have to be run as root. I am working on that bug when it starts annoying too much

```sh
sudo ./pipwm -h
```
prints out commands.  The verbose option "-v" allows to see how bad the quantization error is.

## Examples ##

Controlling servo to zero, min and max
```sh
sudo ./pipwm -v -sa 0
sudo ./pipwm -v -sa -90
sudo ./pipwm -v -sa 90
```

Options -sc sets center, -ss and swing.. and -tr length of period. By default servo mode uses 20000

Specify pulses with lo and hi times
```sh
sudo ./pipwm -lo 80000 -hi 50000 -v
```

Calculate by using hi time and percent
```sh
sudo ./pipwm -hi 9000 -p 5 -v
```

