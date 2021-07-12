package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("127.0.0.1", "3000")
	if err != nil {
		fmt.Println("Failed to listen")
		os.Exit(0)
	}
	defer l.Close()

}
