package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/pedalboard/somb/internal/sushi_rpc"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const midi_port = "pedalboard-midi"

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial("pedalboard-dev:51051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dail sushi grpc port: %v", err)
	}
	defer conn.Close()
	err = connect_to_sushi(ctx, conn)
	if err != nil {
		log.Fatalf("failed to connect to sushi: %v", err)
	}

	defer midi.CloseDriver()

	in, err := midi.FindInPort(midi_port)
	if err != nil {
		log.Fatalf("failed opening midi port: %v", err)
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
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
	sigchan := make(chan os.Signal, 10)

	// listen for ctrl+c
	go signal.Notify(sigchan, os.Interrupt)

	// interrupt has happend
	<-sigchan
	stop()
	os.Exit(0)
}

func connect_to_sushi(ctx context.Context, cc grpc.ClientConnInterface) error {
	sc := sushi_rpc.NewSystemControllerClient(cc)
	v, err := sc.GetSushiVersion(ctx, &sushi_rpc.GenericVoidValue{})
	if err != nil {
		return fmt.Errorf("failed to get sushi version: %w", err)
	}
	log.Printf("connected to sushi version: %v", v.Value)
	return nil
}
