package main

import (
	"fmt"
	"net"
)

const (
	IP   = "127.0.0.1"
	Port = "3000"
)

func main() {
	d, err := net.Dial(IP, Port)
	if err != nil {
		fmt.Println("Failed to connect server")
	}
	

}
