package main

import (
	"flag"
	"github.com/rs/zerolog/log"
    pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
    "github.com/dsrvlabs/vatz/sdk"
    "golang.org/x/net/context"
    "google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
    defaultAddr = "127.0.0.1"
    defaultPort = 9099
    pluginName = "vatz-plugin-vatz-health-checker"
)

var (
	addr string
	port int
	ctx context.Context
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "Listening address")
	flag.IntVar(&port, "port", defaultPort, "Listening port")

	flag.Parse()
}

func main() {
    p := sdk.NewPlugin(pluginName)
    p.Register(pluginFeature)


    ctx = context.Background()

	healthCheck("localhost:9090", ctx)
    if err := p.Start(ctx, addr, port); err != nil {
        log.Info().Str("module", "plugin").Msg("exit")
    }
}

func healthCheck(address string, ctx context.Context) {
	log.Info().Str("module", "plugin").Msgf("dial to %v", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Error().Str("module", "plugin").Msg("dial error")
	}
	log.Info().Str("module", "plugin").Msg(conn.GetState().String())

	healthClient := healthpb.NewHealthClient(conn)
	response, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		log.Error().Str("module", "plugin").Msg("health check response error")
	}

	log.Info().Str("module", "plugin").Msgf("%v", response)
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
    // TODO: Fill here.
    ret := sdk.CallResponse{
        FuncName:   "vatzHealthChecker",
        Message:    "vatz alive!",
        Severity:   pluginpb.SEVERITY_UNKNOWN,
        State:      pluginpb.STATE_NONE,
        AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
    }

	healthCheck("localhost:9090", ctx)
    return ret, nil
}
