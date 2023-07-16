package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/pedalboard/somb/internal/sushi"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const default_midi_port = "pedalboard-midi"
const default_sushi_host_port = "localhost:51051"
const default_channel = 2

func main() {

	var midi_port = flag.String("p", default_midi_port, "midi port")
	var sushi_host_port = flag.String("s", default_sushi_host_port, "sushi host:port")
	var midi_channel = flag.Int("c", default_channel, "midi channel")
	var help = flag.Bool("h", false, "help")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	ctx := context.Background()
	cc, err := grpc.Dial(*sushi_host_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dail sushi grpc port: %v", err)
	}
	defer cc.Close()

	sushi := sushi.NewSushi(cc)
	err = sushi.CheckConnection(ctx)
	if err != nil {
		log.Fatalf("failed to connect to sushi: %v", err)
	}

	defer midi.CloseDriver()
	in, err := midi.FindInPort(*midi_port)
	if err != nil {
		log.Fatalf("failed opening midi port: %v", err)
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch, ctrl, value uint8
		switch {
		case msg.GetControlChange(&ch, &ctrl, &value):
			if ch == uint8(*midi_channel) && ctrl < 10 {
				bypassed := value < 64
				_ = sushi.SetProcessorBypassState(ctx, int32(ctrl), bypassed)
			}

		default:
			// ignore
		}
	})

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	sigchan := make(chan os.Signal, 10)

	// listen for ctrl+c
	go signal.Notify(sigchan, os.Interrupt)

	// interrupt has happend
	<-sigchan
	stop()
	os.Exit(0)
}
