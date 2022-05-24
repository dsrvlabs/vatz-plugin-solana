package main

import (
	"os/exec"
	"bytes"
	"strings"
	"fmt"
	"log"
	"github.com/go-resty/resty/v2"
	"github.com/dsrvlabs/vatz/sdk"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	addr = "0.0.0.0"
	port = 9092
)

type MessageBody struct {
	Jsonrpc	string		`json:"jsonrpc"`
	Id	int		`json:"id"`
	Method	string		`json:"method"`
	Params	[]string	`json:"params"`
}

type Response struct {
	Jsonrpc string		`json:"jsonrpc"`
	Results	Result		`json:"result"`
	Id	int		`json:"id"`
}

type Result struct {
	Context	Context		`json:"context"`
	Value	int		`json:"value"`
}

type Context struct {
	Slot	int		`json:"slot"`
}

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

	pubkey := getSolanaPubkey()

	data := MessageBody{
		Jsonrpc:	"2.0",
		Id:		1,
		Method:		"getBalance",
		Params:		[]string{pubkey},
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		SetResult(&Response{}).
		Post("http://localhost:8899")

	if err != nil {
		log.Fatalf("failed to post message: %v", err)
		return err
	}

	log.Println(pubkey, "Balance: ", resp.Result().(*Response).Results.Value)

	return nil
}

func getSolanaPubkey() string {
	cmd := exec.Command("solana", "address")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to run solana cli: %v", err)
	}

	return strings.ReplaceAll(out.String(), "\n", "")
}
