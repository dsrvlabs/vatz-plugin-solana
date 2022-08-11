module.exports = {
  apps : [{
    name   : "balance_monitor",
    script : "./plugins/balance_monitor/run.sh",
  },{
    name   : "cpu_monitor",
    script : "./plugins/cpu_monitor/run.sh",
  },{
    name   : "disk_monitor",
    script : "./plugins/disk_monitor/run.sh",
  },{
    name   : "health_monitor",
    script : "./plugins/health_monitor/run.sh",
  },{
    name   : "mem_monitor",
    script : "./plugins/mem_monitor/run.sh",
  },{
    name   : "vatz_health_checker",
    script : "./plugins/vatz_health_checker/run.sh",
    args   : "--config ./plugins/vatz_health_checker/default.yaml"
  }]
}
