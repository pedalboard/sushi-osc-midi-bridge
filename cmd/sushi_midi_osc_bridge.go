package main

import (
	"fmt"
	"os"

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

	_, err = midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch, ctrl, vel uint8
		switch {
		case msg.GetControlChange(&ch, &ctrl, &vel):
			fmt.Printf("starting control change %v on channel %v with value %v\n", ctrl, ch, vel)
		default:
			// ignore
		}
	})

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	os.Exit(0)
}
