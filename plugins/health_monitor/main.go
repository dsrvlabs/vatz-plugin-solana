package main

import (
	"fmt"
	"log"
	"github.com/go-resty/resty/v2"
	"github.com/dsrvlabs/vatz/sdk"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	addr = "0.0.0.0"
	port = 9091
)

func main() {
	p := sdk.NewPlugin()
	p.Register(pluginFeature)

	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		fmt.Println("exit")
	}
}

func pluginFeature(info, opt map[string]*structpb.Value) error {
	// TODO: Fill here.

	client := resty.New()
	data := fmt.Sprint("http://127.0.0.1:8899/health")

	resp, err := client.R().Get(data)
	if err != nil {
		log.Fatalf("failed to get response: %v", err)
		return err
	}
	log.Println("Response Info: ", resp.String())

	return nil
}
