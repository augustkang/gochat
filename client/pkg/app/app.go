package app

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
)

type App struct {
	ChatReader *bufio.Reader
	ChatWriter *bufio.Writer
}

func NewApp(conn net.Conn) *App {
	return &App{
		ChatReader: bufio.NewReader(os.Stdin),
		ChatWriter: bufio.NewWriter(conn),
	}
}

func (a *App) Run(conn net.Conn) {
	color.Cyan("Welcome to gochat :D")

	fmt.Print("Please enter your name : ")
	err := a.WriteToConn(a.ReadInput())
	if err != nil {
		fmt.Println("failed to get user name")
	}

	fmt.Println("Please type your command")

	recvChan := make(chan string)
	for {
		go ReadFromServer(conn, recvChan)

		input := a.ReadInput()
		err := a.WriteToConn(input)
		if err != nil {
			os.Exit(1)
		}
	}
}

func (a *App) ReadInput() string {
	input, err := a.ChatReader.ReadString('\n')
	if err != nil {
		color.Red("Failed to get input")
		fmt.Println(err)
		os.Exit(0)
	}
	return input
}

func (a *App) WriteToConn(output string) error {
	_, err := a.ChatWriter.Write([]byte(output))
	if err != nil {
		color.Red("failed to write msg to buffer")
		fmt.Println(err)
		return err
	}
	err = a.ChatWriter.Flush()
	if err != nil {
		color.Red("Failed to send buffer to server")
		fmt.Println(err)
		return err
	}
	return nil
}

func ReadFromServer(conn net.Conn, ch chan string) {
	for {
		select {
		case m := <-ch:
			fmt.Println("Other sent!", m)
		default:

			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("failed to read input from server")
			}
			msg = strings.TrimSpace(msg)
			if err != nil {
				fmt.Println("Failed to get message from server", err)
				os.Exit(1)
			}
			ch <- msg
		}
	}
}
