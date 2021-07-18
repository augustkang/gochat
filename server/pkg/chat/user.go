package chat

import (
	"net"
)

type User struct {
	UserName string
	RoomName string
	Conn     net.Conn
}

func NewUser(name string, conn net.Conn) *User {
	return &User{
		UserName: name,
		Conn:     conn,
	}
}

func (u *User) JoinRoom(name string) {
	u.RoomName = name
}

func (u *User) CreateRoom(name string) string {
	r := NewRoom(name)
	return r.RoomName
}
