package controllers

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http server address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	// The actual websocket connection.
	conn *websocket.Conn
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn)

	fmt.Println("New Client joined the hub!")
	fmt.Println(client)
}
