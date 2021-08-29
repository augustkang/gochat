package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/augustkang/gochat/client/pkg/chatapp"
	"github.com/augustkang/gochat/client/pkg/chatui"
)

func main() {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		fmt.Println("Failed to connect server : ", err)
		panic(err)
	}
	defer conn.Close()

	app := chatapp.NewApp(conn)
	app.InitialPrompt()
	ui, rbox := chatui.GetUI(app.UserName, app)
	var wg sync.WaitGroup
	wg.Add(1)
	go app.Run(conn, ui, rbox, wg)
	if err := ui.Run(); err != nil {
		panic(err)
	}
	wg.Wait()
}
