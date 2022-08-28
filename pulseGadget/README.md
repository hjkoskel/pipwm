# Pulse Gadget

This example program shows how to use pipwm (govattu library) features in real hardware implementation.
Uses SSD1306 (128x64 oled I2C display) and matrix keyboard



## simulated display
If ssd1306 is not available.
It is also possible run this software with
*-simdisp=true* command line option. In that case, updated content of display is printed to stdout.

and/or

## simulated keyboard
If matrix keyboard is not available then it is possible to activate
*-simkey=true* command line option

| Keyboard                       | function in sw |
|--------------------------------|----------------|
| KeyArrowUp                     | UP             |
| KeyArrowDown                   | DOWN           |
| KeyArrowLeft or KeyEsc         | BACK           |
| KeyArrowRight or term.KeyEnter | OK             |
| KeyCtrlC                       | exit           |
| Spacebar                       | ONOFF          |
| 0..9                           | 0..9           |
| , or .                         | decimalpoint   |

## simulate all?

Two kind of simulations

run sofware with command

```sh
./pulseGadget -simhw
```
for getting "total simulation". No pwm signal is generated

```sh
sudo ./pulseGadget -simhw=true -simkey=true
```
for producing signal but with simulated 4x4 keyboard and ssd1306




### 4x4 button matrix ###
Software supports scanned 4x4 matrix keypad connected to GPIO pins


Keys are
|          | **col3** (BCM6) | **col2** (BCM13) | **col1** (BCM19) | **col0** (BCM26) |
|----------|----------|----------|----------|----------|
| **row0** (BCM 12) |    7     |     8    |    9     |    UP    |
| **row1** (BCM 16) |    4     |     5    |    6     |    DOWN    |
| **row2** (BCM 20) |    1     |     2    |    3     |    BACK    |
| **row3** (BCM 21)|    ON    |     0    |    OFF   |    OK    |
