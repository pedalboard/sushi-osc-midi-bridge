package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	defer midi.CloseDriver()

	in, err := midi.FindInPort("pedalboard-midi")
	if err != nil {
		fmt.Println("can't find port")
		return
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch, ctrl, vel uint8
		switch {
		case msg.GetControlChange(&ch, &ctrl, &vel):
			fmt.Printf("starting control change %v on channel %v with value %v\n", ctrl, ch, vel)
		default:
			// ignore
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	time.Sleep(time.Second * 5)

	stop()
}
