package main

import (
    "fmt"
    "github.com/rs/zerolog/log"

    pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
    "github.com/dsrvlabs/vatz/sdk"
    "golang.org/x/net/context"
    "google.golang.org/protobuf/types/known/structpb"
    "github.com/shirou/gopsutil/v3/mem"
)

const (
    addr = "0.0.0.0"
    port = 9095
    pluginName = "vatz-plugin-solana-mem-monitor"
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
        FuncName:   "mem_monitor",
        Message:    "Memory usage warning!",
        Severity:   pluginpb.SEVERITY_UNKNOWN,
        State:      pluginpb.STATE_NONE,
        AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
    }

    v, _ := mem.VirtualMemory()
    log.Info().Str("module", "plugin").Float64("Memory Usage:", v.UsedPercent).Msg("mem_monitor")

    return ret, nil
}
