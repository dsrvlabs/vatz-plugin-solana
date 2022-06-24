package main

import (
    "fmt"
    "log"

    pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
    "github.com/dsrvlabs/vatz/sdk"
    "golang.org/x/net/context"
    "google.golang.org/protobuf/types/known/structpb"
    "github.com/shirou/gopsutil/v3/cpu"
)

const (
    addr = "0.0.0.0"
    port = 9094
    pluginName = "vatz-plugin-solana-cpu-monitor"
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
        FuncName:   "cpu_monitor",
        Message:    "CPU usage warning!",
        Severity:   pluginpb.SEVERITY_UNKNOWN,
        State:      pluginpb.STATE_NONE,
        AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
    }

    util, _ := cpu.Percent(0, false)
    log.Println("CPU Usage: ", util)

    return ret, nil
}
