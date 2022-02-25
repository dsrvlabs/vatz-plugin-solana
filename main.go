package main

import (
        "log"
	plugin_grpc "vatz-plugin-solana/grpc"
)

const (
        servName = "vatz-plugin-solana"
)

func main() {
        log.Println("Start Server: ", servName)

	plugin_grpc.StartServer()

}
