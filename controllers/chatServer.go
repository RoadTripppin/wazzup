package controllers

import (
	"github.com/RoadTripppin/wazzup/models"
	"github.com/google/uuid"
)

type WsServer struct {
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	rooms          map[*Room]bool
	users          []models.Users
	roomRepository models.RoomRepository
	userRepository models.UserRepository
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer(roomRepository models.RoomRepository, userRepository models.UserRepository) *WsServer {
	wsServer := &WsServer{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		rooms:          make(map[*Room]bool),
		roomRepository: roomRepository,
		userRepository: userRepository,
	}

	wsServer.users = userRepository.GetAllUsers()

	return wsServer
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case message := <-server.broadcast:
			server.broadcastToClients(message)

		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.userRepository.AddUser(client)

	server.notifyClientJoined(client)
	server.listOnlineClients(client)
	server.clients[client] = true

	server.users = append(server.users, message.Sender)
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
		server.notifyClientLeft(client)

		for i, user := range server.users {
			if user.GetId() == message.Sender.GetId() {
				server.users[i] = server.users[len(server.users)-1]
				server.users = server.users[:len(server.users)-1]
			}
		}

		server.userRepository.RemoveUser(client)
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}
	if foundRoom == nil {
		// Try to run the room from the repository, if it is found.
		foundRoom = server.runRoomFromRepository(name)
	}

	return foundRoom
}

func (server *WsServer) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)

	server.roomRepository.AddRoom(room)

	go room.RunRoom()
	server.rooms[room] = true

	return room
}

func (server *WsServer) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	server.broadcastToClients(message.encode())
}

func (server *WsServer) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	server.broadcastToClients(message.encode())
}

func (server *WsServer) listOnlineClients(client *Client) {
	for _, user := range server.users {
		message := &Message{
			Action: UserJoinedAction,
			Sender: user,
		}
		client.send <- message.encode()
	}
}

func (server *WsServer) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.ID.String() == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}

func (server *WsServer) runRoomFromRepository(name string) *Room {
	var room *Room
	dbRoom := server.roomRepository.FindRoomByName(name)
	if dbRoom != nil {
		room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())
		room.ID, _ = uuid.Parse(dbRoom.GetId())

		go room.RunRoom()
		server.rooms[room] = true
	}

	return room
}
