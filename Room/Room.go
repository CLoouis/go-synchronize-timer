package Room

import (
	"github.com/gorilla/websocket"
	"learn-websocket/Client"
	"learn-websocket/Counter"
)

type (
	Room struct {
		ID string
		clients map[*Client.Client]bool
		commandChannel chan string
		timerChannel chan string
		counter *Counter.Counter
	}
)

func NewRoom(ID string) *Room {
	return &Room{
		ID: ID,
		clients:      make(map[*Client.Client]bool),
		commandChannel:    make(chan string),
		timerChannel: make(chan string),
		counter: Counter.NewCounter(30),
	}
}

func (room *Room) GetRoomCommandChannel() chan string {
	return room.commandChannel
}

func (room *Room) GetRoomTimerChannel() chan string {
	return room.timerChannel
}

func (room *Room) RegisterNewClient(client *Client.Client) {
	room.clients[client] = true
}

func (room *Room) unregisterClient(client *Client.Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastTimerToClient(timeInformation string) {
	for client := range room.clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(timeInformation))
		if err != nil {
			client.Conn.Close()
			delete(room.clients, client)
		}
	}
}

func (room *Room) Run() {
	go func() {
		for {
			select {
			case message := <- room.commandChannel:
				room.counter.HandleCommand(message, room.timerChannel)
			}
		}
	}()

	go func(){
		for {
			select {
			case message := <- room.timerChannel:
				room.broadcastTimerToClient(message)
			}
		}
	}()
}

