package main

//https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets
import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// Define our message object
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"` //string describe event : left== -10 || right === + 10
}
type Player struct {
	Username string `json:"username"`
	x        int64
	y        int64
}

var clients = make(map[*websocket.Conn]*Player) // connected clients
var broadcast = make(chan Message)              // broadcast channel
// Configure the upgrader
var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	//no cors
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	player := Player{x: 0, y: 0}
	// Register our new client
	clients[ws] = &player
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		newX, _ := strconv.ParseInt(msg.Message, 10, 64)
		clients[ws].x = clients[ws].x + newX
		// Send the newly received message to the broadcast channel
		msg.Message = strconv.FormatInt(int64(clients[ws].x), 10)
		broadcast <- msg
	}
}
func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
func main() {
	fmt.Println("websockets project")
	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	// Start the server on localhost port 8080 and log any errors
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
