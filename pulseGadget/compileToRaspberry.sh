env GOOS=linux GOARCH=arm GOARM=5 go build
scp pulseGadget pi@192.168.1.112:/home/pi
