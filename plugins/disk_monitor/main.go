package main

import (
    "fmt"
    "github.com/rs/zerolog/log"
    "flag"

    pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
    "github.com/dsrvlabs/vatz/sdk"
    "golang.org/x/net/context"
    "google.golang.org/protobuf/types/known/structpb"
    "github.com/shirou/gopsutil/v3/disk"
)

const (
    defaultAddr = "127.0.0.1"
    defaultPort = 9096
    pluginName = "vatz-plugin-solana-disk-monitor"
    defaultUrgent = 95
    defaultWarning = 90
)

var (
    urgent int
    warning int
    addr string
    port int
)

func init() {
    flag.StringVar(&addr, "addr", defaultAddr, "Listening address")
    flag.IntVar(&port, "port", defaultPort, "Listening port")
    flag.IntVar(&urgent, "urgent", defaultUrgent, "Disk Usage Alert threshold")
    flag.IntVar(&warning, "warning", defaultWarning, "Disk Usage Warning threshold")

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
        FuncName:   "getDISKUsage",
        Message:    "Disk usage warning!",
        Severity:   pluginpb.SEVERITY_UNKNOWN,
        State:      pluginpb.STATE_NONE,
        AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
    }

    root, _ := disk.Usage("/")
    log.Debug().Str("module", "plugin").Int("Disk Usage (/):", int(root.UsedPercent)).Int("Urgent", urgent).Int("Warning", warning).Msg("disk_monitor")

    mnt, _ := disk.Usage("/mnt/solana")
    log.Debug().Str("module", "plugin").Int("Disk Usage (/mnt/solana):", int(mnt.UsedPercent)).Int("Urgent", urgent).Int("Warning", warning).Msg("disk_monitor")
    // TODO: grep mountinfo and check disk usage for each mount point

    if int(root.UsedPercent) > urgent {
        var message string
        message = fmt.Sprint("Current Disk Usage (/): ", int(root.UsedPercent), "%, over urgent threshold ", urgent, "%")
        ret = sdk.CallResponse{
            Message:	message,
            Severity:	pluginpb.SEVERITY_CRITICAL,
        }

        log.Warn().Str("module", "plugin").Msg(message)
    }

    if int(mnt.UsedPercent) > urgent {
        var message string
        message = fmt.Sprint("Current Disk Usage (/mnt/solana): ", int(mnt.UsedPercent), "%, over urgent threshold ", urgent, "%")
        ret = sdk.CallResponse{
            Message:	message,
            Severity:	pluginpb.SEVERITY_CRITICAL,
        }

        log.Warn().Str("module", "plugin").Msg(message)
    }

    return ret, nil
}
