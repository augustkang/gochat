package main

import (
	"fmt"
	"net"
	"os"

	"github.com/augustkang/gochat/client/pkg/chatapp"
	"github.com/augustkang/gochat/client/pkg/chatui"
)

func main() {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		fmt.Println("Failed to connect server : ", err)
		os.Exit(0)
	}
	defer conn.Close()

	app := chatapp.NewApp(conn)
	app.InitialPrompt()
	ui, rbox := chatui.GetUI(app.UserName, app)
	go app.Run(conn, ui, rbox)
	if err := ui.Run(); err != nil {
		panic(err)
	}

}
