package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/augustkang/gochat/server/pkg/chat"
	"github.com/fatih/color"
)

const (
	Protocol = "tcp"
	Port     = ":3000"
)

type Server struct {
	Port     string
	RoomList map[string]*chat.Room
	UserList map[string]*chat.User
}

func NewServer(p string) *Server {
	return &Server{
		Port:     p,
		RoomList: make(map[string]*chat.Room),
		UserList: make(map[string]*chat.User),
	}
}

func (s *Server) SendToUser(m string, u *chat.User) {
	w := bufio.NewWriter(u.Conn)
	n, err := w.Write([]byte(m))
	if err != nil {
		fmt.Println(n, err)
	}
	err = w.Flush()
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) Broadcast(m string, u *chat.User) {
	for _, ru := range s.RoomList[u.RoomName].Users {
		if ru.UserName == u.UserName {
			continue
		} else {
			fmt.Println("send message to", ru.UserName)
			go s.SendToUser(m, ru)
		}
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	r, err := reader.ReadString('\n')
	if err != nil {
		color.Red("Failed to get user input. err : ", err)
		os.Exit(1)
	}

	r = strings.TrimSpace(r)
	u := chat.NewUser(r, conn)
	s.UserList[u.UserName] = u
	fmt.Printf("New user %s joined\n", u.UserName)

	for {
		r, err := reader.ReadString('\n')
		if err != nil {
			color.Red("failed to read string from client")
			fmt.Println(err)
			os.Exit(1)
		}

		r = strings.TrimSpace(r)
		input := strings.Split(r, " ")

		switch input[0] {
		case "join":
			exist := s.SearchRoom(input[1])
			if exist {
				fmt.Printf("User %s joined %s\n", u.UserName, input[1])
				r := s.RoomList[input[1]]
				r.Users = append(r.Users, u)
				s.RoomList[input[1]] = r
				u.JoinRoom(input[1])
			} else {
				fmt.Printf("Room name %s doesn't exist. Create and join\n", input[1])
				r := chat.NewRoom(input[1])
				r.Users = append(r.Users, u)
				u.JoinRoom(input[1])
				s.RoomList[input[1]] = r
			}
		case "list":
			if len(s.RoomList) == 0 {
				fmt.Println("There is no room.")
			} else {
				for _, r := range s.RoomList {
					fmt.Println("Room : ", r.RoomName)
				}
			}
		case "user":
			for _, u := range s.UserList {
				fmt.Println("responding user list")
				fmt.Println("user : ", u.UserName)
			}
		case "quit":
			fmt.Printf("User %s exited gochat\n", input[0])
		case "send":
			input[1] = input[1] + "\n"
			s.Broadcast(input[1], u)
		}
	}
}

func (s *Server) SearchRoom(roomName string) bool {
	for _, r := range s.RoomList {
		if r.RoomName == roomName {
			return true
		}
	}
	return false
}

func (s *Server) Run() {
	l, err := net.Listen(Protocol, s.Port)
	if err != nil {
		color.Red("Failed to listen port : ", s.Port)
		os.Exit(0)
	}
	defer l.Close()
	color.Cyan("Server started and listening ")

	for {
		conn, err := l.Accept()

		if err != nil {
			color.Red("Failed to accept client")
			fmt.Println(err)
			os.Exit(0)
		}
		go s.HandleConnection(conn)
	}
}

func main() {
	s := NewServer(Port)
	s.Run()
}
