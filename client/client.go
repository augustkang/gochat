package main

import (
	"fmt"
	"net"
	"os"

	"github.com/augustkang/gochat/client/pkg/app"
)

func main() {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		fmt.Println("Failed to connect server : ", err)
		os.Exit(0)
	}
	defer conn.Close()

	a := app.NewApp(conn)

	a.Run(conn)
}
