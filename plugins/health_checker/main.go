package main

import (
	"time"
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/robfig/cron/v3"
    pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
    "github.com/dsrvlabs/vatz/sdk"
    "golang.org/x/net/context"
    "google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type arrayFlags []string


const (
    defaultAddr = "127.0.0.1"
    defaultPort = 9099
    pluginName = "vatz-plugin-vatz-health-checker"
)

var (
	addr string
	port int
	schedule arrayFlags
	ctx context.Context
)

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "Listening address")
	flag.IntVar(&port, "port", defaultPort, "Listening port")
	flag.Var(&schedule, "schedule", "Schedule with cron expression")

	flag.Parse()
}

func main() {
    p := sdk.NewPlugin(pluginName)
    p.Register(pluginFeature)

    ctx = context.Background()

	c := cron.New(cron.WithLocation(time.UTC))
	for i := 0; i < len(schedule); i++ {
		log.Info().Str("module", "plugin").Msgf("%d, %s", i, schedule[i])
		c.AddFunc(schedule[i], func() { healthCheck("localhost:9090", ctx) })
	}

	c.Start()
	//healthCheck("localhost:9090", ctx)
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
