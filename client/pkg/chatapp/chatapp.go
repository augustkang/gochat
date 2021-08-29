package chatapp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/marcusolsson/tui-go"
)

type App struct {
	UserName   string
	ChatWriter *bufio.Writer
}

func NewApp(conn net.Conn) *App {
	return &App{
		ChatWriter: bufio.NewWriter(conn),
	}
}

func (app *App) InitialPrompt() {
	fmt.Print("Please enter your name : ")
	app.UserName = app.ReadInput()
	err := app.WriteToConn(app.UserName)
	if err != nil {
		fmt.Println("failed to get user name", err)
		return
	}

	fmt.Print("Enter room to chat (ex. room1) : ")
	err = app.WriteToConn(app.ReadInput())
	if err != nil {
		fmt.Println("failed to get user name", err)
		return
	}
}

func (app *App) Run(conn net.Conn, ui tui.UI, cbox *tui.Box, wg sync.WaitGroup) {
	defer wg.Done()
	for {
		m, err := app.ReadFromServer(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		ui.Update(func() {
			cbox.Append(tui.NewHBox(
				tui.NewLabel(m),
				tui.NewSpacer(),
			))
		})
	}
}

func (app *App) ReadFromServer(conn net.Conn) (string, error) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("failed to read input from server")
		return "", err
	}
	msg = strings.TrimSpace(msg)
	if err != nil {
		fmt.Println("Failed to trim space of message", err)
		return "", err
	}
	return msg, nil

}

func (app *App) ReadInput() string {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		color.Red("Failed to get input")
		panic(err)
	}
	return input
}

func (app *App) WriteToConn(output string) error {
	_, err := app.ChatWriter.Write([]byte(output))
	if err != nil {
		color.Red("failed to write msg to buffer")
		fmt.Println(err)
		return err
	}
	err = app.ChatWriter.Flush()
	if err != nil {
		color.Red("Failed to send buffer to server")
		fmt.Println(err)
		return err
	}
	return nil
}
