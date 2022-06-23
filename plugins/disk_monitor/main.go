package main

import (
    "fmt"
    "github.com/rs/zerolog/log"

    pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
    "github.com/dsrvlabs/vatz/sdk"
    "golang.org/x/net/context"
    "google.golang.org/protobuf/types/known/structpb"
    "github.com/shirou/gopsutil/v3/disk"
)

const (
    addr = "0.0.0.0"
    port = 9096
    pluginName = "vatz-plugin-solana-disk-monitor"
)

func main() {
    p := sdk.NewPlugin(pluginName)
    p.Register(pluginFeature)

    ctx := context.Background()
    if err := p.Start(ctx, addr, port); err != nil {
        fmt.Println("exit")
    }
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
    // TODO: Fill here.
    ret := sdk.CallResponse{
        FuncName:   "disk_monitor",
        Message:    "Disk usage warning!",
        Severity:   pluginpb.SEVERITY_UNKNOWN,
        State:      pluginpb.STATE_NONE,
        AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
    }

    root, _ := disk.Usage("/")
    log.Info().Str("module", "plugin").Float64("Disk Usage (/):", root.UsedPercent).Msg("disk_monitor")

    mnt, _ := disk.Usage("/mnt/solana")
    log.Info().Str("module", "plugin").Float64("Disk Usage (/mnt/solana):", mnt.UsedPercent).Msg("disk_monitor")
    // TODO: grep mountinfo and check disk usage for each mount point
    return ret, nil
}
