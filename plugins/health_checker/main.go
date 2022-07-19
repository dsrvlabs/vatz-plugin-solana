package main

import (
	"time"
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/robfig/cron/v3"

	"github.com/dsrvlabs/vatz/manager/notification"
	"github.com/dsrvlabs/vatz/manager/config"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type arrayFlags []string


const (
	defaultAddr = "127.0.0.1"
	defaultPort = 9090
	defaultSchedule = "* * * * *"
)

type notificationInfo struct {
	DiscordSecret	string
	PagerDutySecret	string
}

var (
	addr string
	port int
	schedule arrayFlags
	dispatchManager	= notification.GetDispatcher()
	configFile string
)

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	schedule.Set(defaultSchedule)

	flag.StringVar(&addr, "addr", defaultAddr, "remote address")
	flag.IntVar(&port, "port", defaultPort, "remote port")
	flag.Var(&schedule, "schedule", "Schedule with cron expression, use UTC timezone, space separated ex.) \"* * * * *\"")
	flag.StringVar(&configFile, "config", "default.yaml", "config file path")
	flag.Parse()

	config.InitConfig(configFile)

}

func main() {
	cfg := config.GetConfig()
	vatzConfig := cfg.Vatz

	n := notificationInfo{DiscordSecret: vatzConfig.NotificationInfo.DiscordSecret}
	raddr := fmt.Sprint(addr, ":",  port)
	c := cron.New(cron.WithLocation(time.UTC))

	for i := 0; i < len(schedule); i++ {
		log.Debug().Str("module", "plugin").Msgf("Add cronjob %d, %s", i, schedule[i])
		c.AddFunc(schedule[i], func() { healthCheck(raddr, n.DiscordSecret, context.Background()) })
	}

	c.Start()

	for {
	}
}

func healthCheck(address string, webhook string, ctx context.Context) {
	jsonMessage := notification.ReqMsg{
		FuncName:	"is_vatz_up",
		State:		notification.Success,
		Msg:		"vatz is Alive!!",
		Severity:	notification.Info,
		ResourceType:	"vatz_health_checker",
	}

	log.Info().Str("module", "plugin").Msgf("dial to %v", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Error().Str("module", "plugin").Msg("dial error")
	}

	healthClient := healthpb.NewHealthClient(conn)
	response, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{})

	if err != nil {
		log.Error().Str("module", "plugin").Msg("health check response error")
		jsonMessage = notification.ReqMsg{
			State:		notification.Faliure,
			Msg:		"vatz is down !!",
			Severity:	notification.Critical,
		}

		dispatchManager.SendDiscord(jsonMessage, webhook)
		log.Info().Str("module", "plugin").Msgf("%v", response)

		return
	}

	dispatchManager.SendDiscord(jsonMessage, webhook)

	log.Info().Str("module", "plugin").Msgf("%v", response)
}
