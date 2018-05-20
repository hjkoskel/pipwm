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

Controllin servo to zero, min and max
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

## Hardware support ##

###

### 4x4 button matrix ###
For standalone gadget usage, software supports scanned 4x4 matrix keypad connected to GPIO pins

[]()  |
------|------
Row 1 | row 2


Keys are

[]()    |       |         |
--------|-------|---------|-------
   7    |   8   |   9     |   UP
   4    |   5   |   6     |   DN
   1    |   2   |   3     |   BK
   ON   |   0   |   OFF   |   OK

### GPIO-Button: On/Off toggle ###
At start raspberry checks "not pressed state". Code will decide is switch
normal up or normal false.

Change input and it pulse settings will toggle

Generator starts from off state.

### Button: Enable ###
Single button. Pull pin to gnd = enable, leave open disable



## PulseGadget
This is more complicated example. Uses 128x64 oled I2C display and 4x4 matrix keyboard

## Simulated Pulse Gadget
Example how to develop oled and button interfaces on pc
