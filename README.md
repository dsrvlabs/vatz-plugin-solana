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
```
 # go build
```

4. How to run
- use [pm2](https://pm2.keymetrics.io/docs/usage/quick-start/) to run and manage multiple plugins
- add plugin info into [ecosystem.config.js](./ecosystem.config.js)
```
 # pm2 start ecosystem.config.js
[PM2][WARN] Applications balance_monitor, cpu_monitor, disk_monitor, health_monitor, mem_monitor, vatz_health_checker not running, starting...
[PM2] App [balance_monitor] launched (1 instances)
[PM2] App [cpu_monitor] launched (1 instances)
[PM2] App [disk_monitor] launched (1 instances)
[PM2] App [health_monitor] launched (1 instances)
[PM2] App [mem_monitor] launched (1 instances)
[PM2] App [vatz_health_checker] launched (1 instances)
┌─────┬────────────────────────┬─────────────┬─────────┬─────────┬──────────┬────────┬──────┬───────────┬──────────┬──────────┬──────────┬──────────┐
│ id  │ name                   │ namespace   │ version │ mode    │ pid      │ uptime │ ↺    │ status    │ cpu      │ mem      │ user     │ watching │
├─────┼────────────────────────┼─────────────┼─────────┼─────────┼──────────┼────────┼──────┼───────────┼──────────┼──────────┼──────────┼──────────┤
│ 0   │ balance_monitor        │ default     │ N/A     │ fork    │ 3643402  │ 0s     │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 1   │ cpu_monitor            │ default     │ N/A     │ fork    │ 3643406  │ 0s     │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 2   │ disk_monitor           │ default     │ N/A     │ fork    │ 3643418  │ 0s     │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 3   │ health_monitor         │ default     │ N/A     │ fork    │ 3643431  │ 0s     │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 4   │ mem_monitor            │ default     │ N/A     │ fork    │ 3643442  │ 0s     │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 5   │ vatz_health_checker    │ default     │ N/A     │ fork    │ 3643454  │ 0s     │ 0    │ online    │ 0%       │ 3.2mb    │ root     │ disabled │
└─────┴────────────────────────┴─────────────┴─────────┴─────────┴──────────┴────────┴──────┴───────────┴──────────┴──────────┴──────────┴──────────┘
 # pm2 status
┌─────┬────────────────────────┬─────────────┬─────────┬─────────┬──────────┬────────┬──────┬───────────┬──────────┬──────────┬──────────┬──────────┐
│ id  │ name                   │ namespace   │ version │ mode    │ pid      │ uptime │ ↺    │ status    │ cpu      │ mem      │ user     │ watching │
├─────┼────────────────────────┼─────────────┼─────────┼─────────┼──────────┼────────┼──────┼───────────┼──────────┼──────────┼──────────┼──────────┤
│ 0   │ balance_monitor        │ default     │ N/A     │ fork    │ 3643402  │ 22m    │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 1   │ cpu_monitor            │ default     │ N/A     │ fork    │ 3643406  │ 22m    │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 2   │ disk_monitor           │ default     │ N/A     │ fork    │ 3643418  │ 22m    │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 3   │ health_monitor         │ default     │ N/A     │ fork    │ 3643431  │ 22m    │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 4   │ mem_monitor            │ default     │ N/A     │ fork    │ 3643442  │ 22m    │ 0    │ online    │ 0%       │ 3.1mb    │ root     │ disabled │
│ 5   │ vatz_health_checker    │ default     │ N/A     │ fork    │ 3643454  │ 22m    │ 0    │ online    │ 0%       │ 3.2mb    │ root     │ disabled │
└─────┴────────────────────────┴─────────────┴─────────┴─────────┴──────────┴────────┴──────┴───────────┴──────────┴──────────┴──────────┴──────────┘
 # pm2 stop all
[PM2] Applying action stopProcessId on app [all](ids: [ 0, 1, 2, 3, 4, 5 ])
[PM2] [balance_monitor](0) ✓
[PM2] [cpu_monitor](1) ✓
[PM2] [disk_monitor](2) ✓
[PM2] [health_monitor](3) ✓
[PM2] [mem_monitor](4) ✓
[PM2] [vatz_health_checker](5) ✓
┌─────┬────────────────────────┬─────────────┬─────────┬─────────┬──────────┬────────┬──────┬───────────┬──────────┬──────────┬──────────┬──────────┐
│ id  │ name                   │ namespace   │ version │ mode    │ pid      │ uptime │ ↺    │ status    │ cpu      │ mem      │ user     │ watching │
├─────┼────────────────────────┼─────────────┼─────────┼─────────┼──────────┼────────┼──────┼───────────┼──────────┼──────────┼──────────┼──────────┤
│ 0   │ balance_monitor        │ default     │ N/A     │ fork    │ 0        │ 0      │ 0    │ stopped   │ 0%       │ 0b       │ root     │ disabled │
│ 1   │ cpu_monitor            │ default     │ N/A     │ fork    │ 0        │ 0      │ 0    │ stopped   │ 0%       │ 0b       │ root     │ disabled │
│ 2   │ disk_monitor           │ default     │ N/A     │ fork    │ 0        │ 0      │ 0    │ stopped   │ 0%       │ 0b       │ root     │ disabled │
│ 3   │ health_monitor         │ default     │ N/A     │ fork    │ 0        │ 0      │ 0    │ stopped   │ 0%       │ 0b       │ root     │ disabled │
│ 4   │ mem_monitor            │ default     │ N/A     │ fork    │ 0        │ 0      │ 0    │ stopped   │ 0%       │ 0b       │ root     │ disabled │
│ 5   │ vatz_health_checker    │ default     │ N/A     │ fork    │ 0        │ 0      │ 0    │ stopped   │ 0%       │ 0b       │ root     │ disabled │
└─────┴────────────────────────┴─────────────┴─────────┴─────────┴──────────┴────────┴──────┴───────────┴──────────┴──────────┴──────────┴──────────┘
```
