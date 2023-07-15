package sushi

import (
	"context"
	"fmt"
	"log"

	"github.com/pedalboard/somb/internal/sushi_rpc"
	"google.golang.org/grpc"
)

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
