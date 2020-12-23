package main

import (
	"github.com/gorilla/websocket"
	"learn-websocket/Client"
	"learn-websocket/WsServer"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func main() {
	wsServer := WsServer.NewWsServer()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		handleConnections(wsServer, w, r)
	})

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(wsServer *WsServer.WsServer, w http.ResponseWriter, r *http.Request) {
	ID, ok := r.URL.Query()["id"]

	if !ok || len(ID[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	room := wsServer.FindRoomByID(ID[0])
	if room == nil {
		room = wsServer.CreateRoom(ID[0])
	}

	client := Client.NewClient(ws)
	client.SetRoomCommandChannel(room.GetRoomCommandChannel())

	room.RegisterNewClient(client)
	client.Read()
}