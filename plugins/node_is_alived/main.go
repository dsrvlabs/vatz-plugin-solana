package main

import (
	"fmt"
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/go-resty/resty/v2"
	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	defaultAddr = "127.0.0.1"
	defaultPort = 10002
	pluginName = "vatz-plugin-solana-health-monitor"
)

var (
	addr string
	port int
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "Listening address")
	flag.IntVar(&port, "port", defaultPort, "Listening port")

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
		Message:	"Node status is Healthy",
		Severity:	pluginpb.SEVERITY_UNKNOWN,
		State:		pluginpb.STATE_NONE,
	}

	client := resty.New()
	data := fmt.Sprint("http://127.0.0.1:8899/health")

	resp, err := client.R().Get(data)
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("failed to get response: %v", err)
		ret = sdk.CallResponse{
			Message:	"failed to get response",
			State:		pluginpb.STATE_FAILURE,
		}

		return ret, err
	}

	if resp.String() != "ok" {
		var message string
		message = "Node status is unhealthy " + resp.String()
		ret = sdk.CallResponse{
			Message:	message,
			Severity:	pluginpb.SEVERITY_CRITICAL,
		}

		log.Warn().Str("module", "plugin").Msg(message)
	} else {
		log.Debug().Str("module", "plugin").Msg("Node status is healthy")
	}

	return ret, nil
}
