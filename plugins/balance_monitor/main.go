package main

import (
	"flag"
	"os/exec"
	"bytes"
	"strings"
	"strconv"
	"github.com/rs/zerolog/log"
	"github.com/go-resty/resty/v2"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	defaultAddr = "127.0.0.1"
	defaultPort = 9092
	defaultWarning = 15
	defaultUrgent = 5
	pluginName = "vatz-plugin-solana-balance-monitor"
	lamport = 1000000000
)

var (
	addr string
	port int
	warning int
	urgent int
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

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "Listening address")
	flag.IntVar(&port, "port", defaultPort, "Listening port")
	flag.IntVar(&warning, "warning", defaultWarning, "Warning threshold for insufficient SOL Balance")
	flag.IntVar(&urgent, "urgent", defaultUrgent, "Urgent threshold for insufficient SOL Balance")

	flag.Parse()
}

func main() {
	p := sdk.NewPlugin(pluginName)
	p.Register(pluginFeature)

	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		log.Info().Str("module", "plugin").Msg("exit")
	}
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
	// TODO: Fill here.
	ret := sdk.CallResponse{
		FuncName:   info["execute_method"].GetStringValue(),
		Message:	"SOL Balance for voting fee",
		Severity:	pluginpb.SEVERITY_UNKNOWN,
		State:		pluginpb.STATE_NONE,
	}

	pubkey := getSolanaPubkey()

	data := MessageBody{
		Jsonrpc:	"2.0",
		Id:			1,
		Method:		"getBalance",
		Params:		[]string{pubkey},
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		SetResult(&Response{}).
		Post("http://127.0.0.1:8899")

	if err != nil {
		log.Error().Str("module", "plugin").Msgf("failed to post message: %v", err)
		ret = sdk.CallResponse{
			Message:	"failed to get response",
			State:		pluginpb.STATE_FAILURE,
		}

		return ret, err
	}

	balance := resp.Result().(*Response).Results.Value
	log.Debug().Str("module", "plugin").Msgf("current balance: %d", balance)

	if balance < urgent {
		var message string
		message = "insufficient SOL Balance. Charge Voting fee within a day. current Balance : " + strconv.Itoa(balance)

		ret = sdk.CallResponse{
			Message:	message,
			Severity:	pluginpb.SEVERITY_CRITICAL,
		}
		log.Error().Str("module", "plugin").Msg(message)
	} else if balance < warning {
		var message string
		message = "insufficient SOL Balance. Charge Voting fee within a week. current Balance : " + strconv.Itoa(balance)

		ret = sdk.CallResponse{
			Message:	message,
			Severity:	pluginpb.SEVERITY_WARNING,
		}
		log.Warn().Str("module", "plugin").Msg(message)
	}

	return ret, nil
}

func getSolanaPubkey() string {
	cmd := exec.Command("solana", "address")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("failed to run solana cli: %v", err)
	}

	return strings.ReplaceAll(out.String(), "\n", "")
}
