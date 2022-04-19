package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RoadTripppin/wazzup/config"
	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

//var addr = flag.String("addr", ":8080", "http server address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	// The actual websocket connection.
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
	Name     string    `json:"name"`
	ID       uuid.UUID `json:"id"`
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string, id string) *Client {
	return &Client{
		ID:       uuid.MustParse(id),
		Name:     name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
	}
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		client.handleNewMessage(jsonMessage)
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:

			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	//client.wsServer.unregister <- client
	//for room := range client.rooms {
	//	room.unregister <- client
	//}
	close(client.send)
	client.conn.Close()
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

	// name, ok := r.URL.Query()["name"]
	user_id, ok := r.URL.Query()["ID"]

	if !ok || len(user_id[0]) < 1 {
		log.Println("Url Param 'ID' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(user_id[0])
	// Get User and send required details to newClient
	db := config.InitDB()
	user := &helpers.User{}

	row := db.QueryRow("SELECT id, name FROM user WHERE id = ?", user_id[0])

	if err := row.Scan(&user.Id, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			// panic(err)
			log.Println("No user found")
			return
			// return map[string]interface{}{"message": "User not found"}
		}
		//panic(err)
	}

	client := newClient(conn, wsServer, user.Name, user.Id)

	go client.writePump()
	go client.readPump()

	wsServer.register <- client

	// Get all rooms user is in and send joinRoom message(subscribe user to it)
	fmt.Println("Fetching rooms to auto-join")
	row = db.QueryRow("SELECT rooms from user where id = ?", client.GetId())
	var rooms string
	if err := row.Scan(&rooms); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found")
			return
		}
	}

	if rooms != "" {
		var roomArray []string
		if strings.Contains(rooms, ":;;:") {
			fmt.Println("In split")
			roomArray = strings.Split(rooms, ":;;:")
		} else {
			roomArray = append(roomArray, rooms)
		}

		for _, room := range roomArray {
			// Get room Name from DB with room.ID
			var roomData map[string]interface{}
			json.Unmarshal([]byte(room), &roomData)
			fmt.Println("Room ID: ", roomData["id"])

			row = db.QueryRow("SELECT name FROM room WHERE id = ?", roomData["id"])

			var roomname string
			if err := row.Scan(&roomname); err != nil {
				if err == sql.ErrNoRows {
					log.Println("Error: No room found")
					return
				}
			}

			message := Message{
				Message: roomname,
			}
			fmt.Println("Room Name: ", roomname)

			client.handleJoinRoomMessage(message)
		}
	}
}

func (client *Client) handleNewMessage(jsonMessage []byte) {

	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	// Attach the client object as the sender of the messsage.
	message.Sender = client

	switch message.Action {
	case SendMessageAction:
		// The send-message action, this will send messages to a specific room now.
		// Which room wil depend on the message Target
		roomID := message.Target.GetId()
		// Use the ChatServer method to find the room, and if found, broadcast!
		if room := client.wsServer.findRoomByID(roomID); room != nil {
			// Add timestamp to message, stringify it and then send it
			messageText := map[string]interface{}{
				"text":      message.Message,
				"timestamp": time.Now(),
			}

			messageString, _ := json.Marshal(messageText)
			message.Message = string(messageString)

			room.broadcast <- &message

			// Insert message object to messages column of room table
			client.insertMessageToDB(room, message)

		}
	// We delegate the join and leave actions.
	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)

	case JoinRoomPrivateAction:
		client.handleJoinRoomPrivateMessage(message)
	}
}

func (client *Client) handleJoinRoomMessage(message Message) {
	//Will have to find room from db //
	roomName := message.Message

	client.joinRoom(roomName, nil)
}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.wsServer.findRoomByID(message.Message)
	if room == nil {
		return
	}
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}

func (client *Client) joinRoom(roomName string, sender models.Users) *Room {

	// Get roomname from DB instead of existing functionality
	// fmt.Println("In join room method")
	// db := config.InitDB()
	// room := &Room{}

	// row := db.QueryRow("SELECT id, name, private FROM room WHERE name = ?", roomName)

	// // var newRoom helpers.Room
	// if err := row.Scan(&room.ID, &room.Name, &room.Private); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		// return map[string]interface{}{"message": "User not found"}

	// 		// If no room found, create room and insert into DB as well
	// 		fmt.Println("Room not found")
	// 		room = client.wsServer.createRoom(roomName, sender != nil)
	// 		// Insert room into DB
	// 		db := config.InitDB()
	// 		stmt, err := db.Prepare("INSERT INTO room(id, name, private) values(?,?,?)")
	// 		helpers.CheckErr(err)

	// 		_, err = stmt.Exec(room.ID.String(), room.Name, room.Private)
	// 		helpers.CheckErr(err)

	// 		defer db.Close()
	// 	}
	// 	//panic(err)
	// }
	// fmt.Println(room.ID)
	// fmt.Println(room.Name)
	// if room == nil {
	// room = client.wsServer.createRoom(roomName, sender != nil)
	// }

	// room := client.wsServer.findRoomByName(roomName)
	// if room == nil {
	// 	room = client.wsServer.createRoom(roomName, sender != nil)
	// }

	// Don't allow to join private rooms through public room message
	//fmt.Println("sender", sender.GetName())

	fmt.Println("In join room")
	room := client.wsServer.findRoomByName(roomName)

	//fmt.Println("Room found: ", room.ID)
	// log.Println("In join room method")
	if room == nil {
		room = client.wsServer.createRoom(roomName, sender != nil)
	}

	// fmt.Println("Room Private: ", room.Private)
	// if sender == nil && room.Private {
	// 	fmt.Println("Inside nil condition")
	// 	return nil
	// }

	if !client.isInRoom(room) {
		fmt.Println("Join Room: ", client.Name)
		client.rooms[room] = true
		room.register <- client

		// Update client/user record in db to contain room info
		if sender != nil {
			client.addRoomToDB(sender.GetId(), room)

			client.notifyRoomJoined(room, sender)
		}
	}
	// fmt.Println("return ->", client.Name)
	return room
}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}
	return false
}

func (client *Client) notifyRoomJoined(room *Room, sender models.Users) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}

func (client *Client) GetName() string {
	return client.Name
}

func (client *Client) GetId() string {
	return client.ID.String()
}

func (client *Client) handleJoinRoomPrivateMessage(message Message) {
	// instead of searching for a client, search for User by the given ID.

	// Modify to find User from DB and use that ID..not server array
	// message.Message here is ID as string
	db := config.InitDB()
	user := &helpers.User{}
	fmt.Println("In Join Private Room method: ")
	fmt.Println(message.Message)
	row := db.QueryRow("SELECT id, name, email, password, profilepic FROM user WHERE id = ?", message.Message)

	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.ProfilePic); err != nil {
		if err == sql.ErrNoRows {
			// return map[string]interface{}{"message": "User not found"}
			return
		}
		//panic(err)
	}

	// target := client.wsServer.findUserByID(message.Message)
	// if target == nil {
	// 	return
	// }

	// fmt.Println(user.GetId())
	// if target.GetId() == user.GetId() {
	// 	fmt.Println("ID is true")
	// }

	roomName := message.Message + client.ID.String()

	// joinedRoom := client.joinRoom(roomName, target)
	joinedRoom := client.joinRoom(roomName, user)

	if joinedRoom != nil {
		// client.inviteTargetUser(target, joinedRoom)
		client.inviteTargetUser(user, joinedRoom)

		// *** Update Sender and Target User's record to contain room info in DB here ***
		log.Println(joinedRoom.Name + " : Can update this section now")
	}
}

func (client *Client) inviteTargetUser(target models.Users, room *Room) {
	inviteMessage := &Message{
		Action:  JoinRoomPrivateAction,
		Message: target.GetId(),
		Target:  room,
		Sender:  client,
	}

	if err := config.Redis.Publish(Contex, PubSubGeneralChannel, inviteMessage.encode()).Err(); err != nil {
		log.Println(err)
	}

}

// Use this from join room to insert into DB
func (client *Client) addRoomToDB(target string, room *Room) {
	// Room info: {Room ID, Room Name(Sender Name), Sender Profile Pic}
	db := config.InitDB()
	var rooms string

	roomRow := db.QueryRow("SELECT rooms FROM user WHERE id = ?", client.GetId())
	if err := roomRow.Scan(&rooms); err != nil {
		if err == sql.ErrNoRows {
			// panic(err)
			log.Println("No user found: In Join Room")
			return
			// return map[string]interface{}{"message": "User not found"}
		}
		//panic(err)
	}

	user := helpers.User{}
	targetRow := db.QueryRow("SELECT name, profilepic FROM user WHERE id =?", target)
	if err := targetRow.Scan(&user.Name, &user.ProfilePic); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Target found: In Join Room")
			return
		}
	}

	roomData := map[string]interface{}{
		"id":         room.GetId(),
		"name":       user.Name,
		"profilepic": user.ProfilePic,
	}

	roomString, _ := json.Marshal(roomData)

	if rooms == "" {
		rooms = string(roomString)
	} else {
		rooms = rooms + ":;;:" + string(roomString)
	}

	stmt, err := db.Prepare("UPDATE user SET rooms = ? WHERE id = ?")
	helpers.CheckErr(err)

	_, err = stmt.Exec(rooms, client.GetId())
	helpers.CheckErr(err)
}

func (client *Client) insertMessageToDB(room *Room, message Message) {
	db := config.InitDB()
	var messages string

	roomRow := db.QueryRow("SELECT messages FROM room WHERE id = ?", room.GetId())
	if err := roomRow.Scan(&messages); err != nil {
		if err == sql.ErrNoRows {
			// panic(err)
			log.Println("No room found: while inserting image")
			return
			// return map[string]interface{}{"message": "User not found"}
		}
		//panic(err)
	}

	var messageText map[string]interface{}

	if err := json.Unmarshal([]byte(message.Message), &messageText); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	messageData := map[string]interface{}{
		"text":       messageText["text"],
		"timestamp":  messageText["timestamp"],
		"senderID":   client.GetId(),
		"senderName": client.GetName(),
	}

	messageString, _ := json.Marshal(messageData)

	if messages == "" {
		messages = string(messageString)
	} else {
		messages = messages + ":;;:" + string(messageString)
	}

	stmt, err := db.Prepare("UPDATE room SET messages = ? WHERE id = ?")
	helpers.CheckErr(err)

	_, err = stmt.Exec(messages, room.GetId())
	helpers.CheckErr(err)
}
