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
	r, err = reader.ReadString('\n')
	if err != nil {
		color.Red("Failed to get user input. err : ", err)
		os.Exit(1)
	}
	r = strings.TrimSpace(r)
	joinCmd := strings.Split(r, " ")
	exist := s.SearchRoom(joinCmd[1])
	if exist {
		fmt.Printf("User %s joined %s\n", u.UserName, joinCmd[1])
		r := s.RoomList[joinCmd[1]]
		r.Users = append(r.Users, u)
		s.RoomList[joinCmd[1]] = r
		u.JoinRoom(joinCmd[1])
		msg := "<<User " + u.UserName + " has joined this room! >>\n"
		s.Broadcast(msg, u)
	} else {
		fmt.Printf("Room name %s doesn't exist. Create and join\n", joinCmd[1])
		r := chat.NewRoom(joinCmd[1])
		r.Users = append(r.Users, u)
		u.JoinRoom(joinCmd[1])
		s.RoomList[joinCmd[1]] = r
	}

	for {
		r, err := reader.ReadString('\n')
		if err != nil {
			color.Red("failed to read string from client")
			fmt.Println(err)
			os.Exit(1)
		}
		r = strings.TrimSpace(r)

		r = u.UserName + " : " + r + "\n"
		fmt.Print(r)
		s.Broadcast(r, u)

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
