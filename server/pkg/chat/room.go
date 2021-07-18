package chat

type Room struct {
	RoomName string
	Users    []*User
}

func NewRoom(name string) (r *Room) {
	return &Room{
		RoomName: name,
	}
}
