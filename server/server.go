package main

import (
	"bufio"
	"fmt"
	"net"
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
		panic(err)
	}
	err = w.Flush()
	if err != nil {
		fmt.Println(err)
		panic(err)
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
		panic(err)
	}
	r = strings.TrimSpace(r)
	u := chat.NewUser(r, conn)
	s.UserList[u.UserName] = u
	fmt.Printf("New user %s joined\n", u.UserName)
	r, err = reader.ReadString('\n')
	if err != nil {
		color.Red("Failed to get user input. err : ", err)
		panic(err)
	}
	roomName := strings.TrimSpace(r)
	exist := s.SearchRoom(roomName)
	if exist {
		fmt.Printf("User %s joined %s\n", u.UserName, roomName)
		r := s.RoomList[roomName]
		r.Users = append(r.Users, u)
		s.RoomList[roomName] = r
		u.JoinRoom(roomName)
		msg := "<<User " + u.UserName + " has joined this room! >>\n"
		s.Broadcast(msg, u)
	} else {
		fmt.Printf("Room name %s doesn't exist. Create and join\n", roomName)
		r := chat.NewRoom(roomName)
		r.Users = append(r.Users, u)
		u.JoinRoom(roomName)
		s.RoomList[roomName] = r
	}

	for {
		r, err := reader.ReadString('\n')
		if err != nil {
			color.Red("failed to read string from client")
			fmt.Println(err)
			panic(err)
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
		panic(err)
	}
	defer l.Close()
	color.Cyan("Server started and listening ")

	for {
		conn, err := l.Accept()

		if err != nil {
			color.Red("Failed to accept client")
			fmt.Println(err)
			panic(err)
		}
		go s.HandleConnection(conn)
	}
}

func main() {
	s := NewServer(Port)
	s.Run()
}
