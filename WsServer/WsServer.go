package WsServer

import (
	"learn-websocket/Room"
)

type (
	WsServer struct {
		rooms map[*Room.Room]bool
	}
)

func NewWsServer() *WsServer {
	return &WsServer{rooms:make(map[*Room.Room]bool)}
}

func (server *WsServer) FindRoomByID(ID string) *Room.Room {
	var foundRoom *Room.Room
	for room := range server.rooms {
		if room.ID == ID {
			foundRoom = room
			break
		}
	}
	return foundRoom
}

func (server *WsServer) CreateRoom(ID string) *Room.Room {
	room := Room.NewRoom(ID)
	go room.Run()
	server.rooms[room] = true

	return room
}
