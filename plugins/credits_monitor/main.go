package main

import (
	"fmt"
	"os/exec"
	"bytes"
	"strings"
	"log"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/go-resty/resty/v2"
	"github.com/dsrvlabs/vatz/sdk"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	addr = "0.0.0.0"
	port = 9093
    pluginName = "vatz-plugin-solana-credits-checker"
	votePubkey = "4kARmfWrpgyCABPHoiSxZh2jMf7GeABX7AhMBXVC8Wds"
)

type MessageBody struct {
    Jsonrpc string      `json:"jsonrpc"`
    Id  int     `json:"id"`
    Method  string      `json:"method"`
    Params  Param    `json:"params"`
}

type Param struct {
	Votepubkey	[]string	`json:"votePubkey"`
}

type Response struct {
    Jsonrpc string      `json:"jsonrpc"`
    Results Result      `json:"result"`
    Id  int     `json:"id"`
}

type Result struct {
    Currents Current     `json:"current"`
    Delinquent   []string     `json:"delinquent"`
}

type Current struct {
    Commission    int     `json:"commission"`
	EpochCredits	[][]int	`json:"epochCredits"`
}

func main() {
	p := sdk.NewPlugin(pluginName)
	p.Register(pluginFeature)

	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		fmt.Println("exit")
	}

	log.Println(resp.Result().(*Response).Results.Currents.EpochCredits)
	log.Println("current epoch: ", getSolanaEpoch())

}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
	// TODO: Fill here.

	ret := sdk.CallResponse {
		FuncName:	"credits_monitor",
		Message:	"Credits Warning!",
		Severity:	pluginpb.SEVERITY_WARNING,
		State:		pluginpb.STATE_SUCCESS,
		AlertTypes:	[]pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	param := Param {
		Votepubkey:	[]string{votePubkey},
	}

	data := MessageBody {
		Jsonrpc:	"2.0",
		Id:			1,
		Method:		"getVoteAccounts",
		Params:		param,
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		SetResult(&Response{}).
		Post("http://162.55.243.228:8899")

	if err != nil {
		log.Fatalf("failed to post message: %v", err)
		return ret, err
	}

	log.Println("credits: ", resp.Result())
	log.Println("current epoch: ", getSolanaEpoch())
	return ret, nil
}

func getSolanaEpoch() string {
	cmd := exec.Command("solana", "epoch")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to run solana cli: %v", err)
	}

	return strings.ReplaceAll(out.String(), "\n", "")
}
