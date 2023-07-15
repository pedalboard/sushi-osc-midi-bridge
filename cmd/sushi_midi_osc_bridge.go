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
	conn, err := grpc.Dial("localhost:51051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dail sushi grpc port: %v", err)
	}
	defer conn.Close()

	sushi := NewSushi(conn)

	err = sushi.CheckConnection(ctx)
	if err != nil {
		log.Fatalf("failed to connect to sushi: %v", err)
	}

	defer midi.CloseDriver()

	in, err := midi.FindInPort(midi_port)
	if err != nil {
		log.Fatalf("failed opening midi port: %v", err)
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch, ctrl, value uint8
		switch {
		case msg.GetControlChange(&ch, &ctrl, &value):
			if ch == 2 && ctrl < 10 {
				bypassed := value < 64
				_ = sushi.SetProcessorBypassState(ctx, int32(ch), bypassed)
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

type Sushi struct {
	cc grpc.ClientConnInterface
}

func NewSushi(cc grpc.ClientConnInterface) *Sushi {
	return &Sushi{cc: cc}
}

func (s *Sushi) CheckConnection(ctx context.Context) error {
	sc := sushi_rpc.NewSystemControllerClient(s.cc)
	v, err := sc.GetSushiVersion(ctx, &sushi_rpc.GenericVoidValue{})
	if err != nil {
		return fmt.Errorf("failed to get sushi version: %w", err)
	}
	log.Printf("connected to sushi version: %v", v.Value)
	return nil
}

func (s *Sushi) SetProcessorBypassState(ctx context.Context, id int32, bypassed bool) error {

	agc := sushi_rpc.NewAudioGraphControllerClient(s.cc)

	_, err := agc.SetProcessorBypassState(ctx, &sushi_rpc.ProcessorBypassStateSetRequest{
		Processor: &sushi_rpc.ProcessorIdentifier{Id: id},
		Value:     bypassed,
	})
	if err != nil {
		return fmt.Errorf("failed to SetProcessorBypassState: %w", err)
	}
	log.Printf("SetProcessorBypassState for processor %v to %v", id, bypassed)
	return nil
}
