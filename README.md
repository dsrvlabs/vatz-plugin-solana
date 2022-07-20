# vatz-plugin-solana
Vatz Plugin for Solana.

1. Plugin configuration
```
plugins_infos:
  default_verify_interval:  15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_name: "vatz-plugin-solana-health-monitor"
      plugin_address: "localhost"
      verify_interval: 60
      execute_interval: 60
      plugin_port: 9091
      executable_methods:
        - method_name: "getHealth"
    - plugin_name: "vatz-plugin-solana-balance-monitor"
      plugin_address: "localhost"
      verify_interval: 60
      execute_interval: 60
      plugin_port: 9092
      executable_methods:
        - method_name: "getBalance"
    - plugin_name: "vatz-plugin-solana-cpu-monitor"
      plugin_address: "localhost"
      verify_interval: 60
      execute_interval: 30
      plugin_port: 9094
      executable_methods:
        - method_name: "getCPUUsage"
    - plugin_name: "vatz-plugin-solana-mem-monitor"
      plugin_address: "localhost"
      verify_interval: 60
      execute_interval: 30
      plugin_port: 9095
      executable_methods:
        - method_name: "getMEMUsage"
    - plugin_name: "vatz-plugin-solana-disk-monitor"
      plugin_address: "localhost"
      verify_interval: 60
      execute_interval: 60
      plugin_port: 9096
      executable_methods:
        - method_name: "getDiskUsage"
```

2. plugin directory structure
```
root@solana-testnet-validator-hetzner ~/vatz/vatz-plugin-solana # tree .
.
├── logs
│   ├── balance_monitor.log
│   ├── cpu_monitor.log
│   ├── disk_monitor.log
│   ├── health_monitor.log
│   └── mem_monitor.log
├── plugins
│   ├── balance_monitor
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── main.go
│   │   ├── run.sh
│   │   └── stop.sh
│   ├── cpu_monitor
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── main.go
│   │   ├── run.sh
│   │   └── stop.sh
│   ├── disk_monitor
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── main.go
│   │   ├── run.sh
│   │   └── stop.sh
│   ├── health_monitor
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── main.go
│   │   ├── run.sh
│   │   └── stop.sh
│   └── mem_monitor
│       ├── go.mod
│       ├── go.sum
│       ├── main.go
│       ├── run.sh
│       └── stop.sh
└── README.md
```

3. How to build

4. How to run
