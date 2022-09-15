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
	nodeStatus bool = true
	schedule arrayFlags
	dispatchManager	= notification.GetDispatcher()
	configFile string
	hostname string
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
	flag.Var(&schedule, "schedule", "Schedule with cron expression, use UTC timezone, space separated ex.) \"0 0 * * *\" (default \"* * * * *\")")
	flag.StringVar(&configFile, "config", "default.yaml", "config file path")
	flag.Parse()

	config.InitConfig(configFile)

}

func main() {
	cfg := config.GetConfig()
	vatzConfig := cfg.Vatz

	n := notificationInfo{DiscordSecret: vatzConfig.NotificationInfo.DiscordSecret}
	raddr := fmt.Sprint(addr, ":",  port)

	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatal().Str("module", "plugin").Msgf("Couldn't get hostname: %v", err)
	}

	ctx := context.Background()
	c := cron.New(cron.WithLocation(time.UTC))

	c.AddFunc(schedule[0], func() { heartBeat(raddr, n.DiscordSecret, ctx) })
	log.Debug().Str("module", "plugin").Msgf("Add cronjob 0, %s", schedule[0])

	for i := 1; i < len(schedule); i++ {
		log.Debug().Str("module", "plugin").Msgf("Add cronjob %d, %s", i, schedule[i])
		c.AddFunc(schedule[i], func() { scheduledCheck(raddr, n.DiscordSecret, ctx) })
	}


	c.Start()

	for {
	}
}

func heartBeat(address string, webhook string, ctx context.Context) {
	// used for checking every minute as default
	response, err := healthCheck(address, ctx)
	if err != nil {
		log.Error().Str("module", "plugin").Msg("health check response error")
		message := fmt.Sprint("[", hostname, "]\n", "vatz is down!!")
		jsonMessage := notification.ReqMsg{
			FuncName:	"is_vatz_up",
			State:		notification.Faliure,
			Msg:		message,
			Severity:	notification.Critical,
			ResourceType:	"vatz_health_checker",
		}

		dispatchManager.SendDiscord(jsonMessage, webhook)
		log.Error().Str("module", "plugin").Msgf("%v", response)
		nodeStatus = false
	} else {
		log.Info().Str("module", "plugin").Msgf("%v", response)
		if nodeStatus == false {
			message := fmt.Sprint("[", hostname, "]\n", "vatz is back !!")
			jsonMessage := notification.ReqMsg{
				FuncName:	"is_vatz_up",
				State:		notification.Success,
				Msg:		message,
				Severity:	notification.Info,
				ResourceType:	"vatz_health_checker",
			}

			dispatchManager.SendDiscord(jsonMessage, webhook)
		}
		nodeStatus = true
	}
}

func scheduledCheck(address string, webhook string, ctx context.Context) {
	// used for checking scheduled cronjob except default
	response, err := healthCheck(address, ctx)
	if err == nil {
		log.Info().Str("module", "plugin").Msg("scheduled health check")
		message := fmt.Sprint("[", hostname, "]\n", "vatz is Alive!!")
		jsonMessage := notification.ReqMsg{
			FuncName:	"is_vatz_up",
			State:		notification.Success,
			Msg:		message,
			Severity:	notification.Info,
			ResourceType:	"vatz_health_checker",
		}

		dispatchManager.SendDiscord(jsonMessage, webhook)
		log.Info().Str("module", "plugin").Msgf("%v", response)
	}
}

func healthCheck(address string, ctx context.Context) (*healthpb.HealthCheckResponse, error) {
	log.Info().Str("module", "plugin").Msgf("dial to %v", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Error().Str("module", "plugin").Msg("dial error")
		return nil, err
	}

	healthClient := healthpb.NewHealthClient(conn)

	response, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{})

	return response, err
}
