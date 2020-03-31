package main

import (
	_ "encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Define our message object
type actionMessage struct {
	X string // "0" || "10" ||  "-10"
	Y string // "0" || "10" ||  "-10"
}

type screenWH struct {
	Width int // "0" || "10" ||  "-10"
	Height int // "0" || "10" ||  "-10"
}

type point struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type Player struct {
	Id            uuid.UUID `json"id"`
	ExceptionType string    `json:"exceptionType"` //add id : uuid
	Color         [3]int    `json:"color"`
	Size          int       `json:"size"`
	P             point     `json:"p"`
	Show          bool      `json:"show"`
}

type Exception struct {
	Id            uuid.UUID `json"id"`
	ExceptionType string    `json:"exceptionType"`
	Show          bool      `json:"show"`
	P             point     `json:"p"`
}

var exceptionsMap = struct {
	sync.RWMutex
	m map[point]Exception
}{m: make(map[point]Exception)}

var s = rand.NewSource(time.Now().UnixNano())
var exceptionsTypes = [3]string{"NullPointerException", "DivideByZeroException", "IOException"}
var clients = make(map[*websocket.Conn]*Player) // connected clients
var broadcastPlayers = make(chan Player)
var broadcastException = make(chan Exception)
var upgrader = websocket.Upgrader{}

func handleNewPlayer(ws *websocket.Conn) {
	r2 := rand.New(s)
	player := Player{Id: uuid.New(), P: point{X: int64(r2.Intn(300)), Y: 0}, Size: 50, Show: true, ExceptionType: exceptionsTypes[rand.Intn(3)], Color: [3]int{r2.Intn(256), r2.Intn(256), r2.Intn(256)}}


	//send to new player all current players


	fmt.Println("new player")
	fmt.Println(player)

	ws.WriteJSON(player)
	
	for key := range clients {
		tempP := *clients[key]
		fmt.Println(tempP)
		err := ws.WriteJSON(tempP)
		if err != nil {
			log.Printf( "73 error: %v", err)
			ws.Close()
			delete(clients, ws)
		}
	}

	//TODO send to new player all current exceptions

	clients[ws] = &player
	broadcastPlayers <- player //broadcast new player



	var msg screenWH
	// Read in a new message as JSON and map it to a screenWH object
	err := ws.ReadJSON(&msg)
	if err != nil {
		log.Printf(" 102 error: %v", err)
		var plyrMsg = clients[ws]
		plyrMsg.Show = false
		broadcastPlayers <- *plyrMsg
		delete(clients, ws)
	}
	fmt.Println("received msg:")
	fmt.Println(msg)


}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //no cors

	ws, err := upgrader.Upgrade(w, r, nil) // Upgrade initial GET request to a websocket
	if err != nil {
		log.Fatal(err)
	}

	handleNewPlayer(ws)

	defer ws.Close() // Make sure we close the connection when the function returns

	for {
		var msg actionMessage
		// Read in a new message as JSON and map it to a incomingMessage object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf(" 102 error: %v", err)
			var plyrMsg = clients[ws]
			plyrMsg.Show = false
			broadcastPlayers <- *plyrMsg
			delete(clients, ws)
			break
		}
		fmt.Println("received msg:")
		fmt.Println(msg)

		fmt.Println("player curr values")
		fmt.Println(clients[ws])
		//INSERT HERE  call to check collision function
		newX, _ := strconv.ParseInt(msg.X, 10, 64)
		clients[ws].P.X = int64(clients[ws].P.X) + newX

			newY, _ :=  strconv.ParseInt(msg.Y, 10, 64)
			clients[ws].P.Y = int64(clients[ws].P.Y) + newY
		fmt.Println("player with new values ")
		fmt.Println(*clients[ws])
		broadcastPlayers <- *clients[ws]
	}
}
func broadcastMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcastPlayers

		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("134 error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}

	}
}

func exceptiosMapHandler() {
	//Thread.
	// every 5 sec:  create new ex ->add to eXarr -> broadcast to users
	//exery 10 sec:  choose rand ex ->remove from exArr ->broadcast to users

	var r = rand.New(s)
	addExTicker := time.NewTicker(5 * time.Second)
	go func() {
		for t := range addExTicker.C {
			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
			var newEx = Exception{Id: uuid.New(), ExceptionType: exceptionsTypes[r.Intn(3)], P: point{X: int64(r.Intn(255)), Y: int64(r.Intn(255))}, Show: true}
			exceptionsMap.Lock() //take the write lock
			exceptionsMap.m[newEx.P] = newEx
			exceptionsMap.Unlock() // release the write lock
			fmt.Println(newEx)
			broadcastException <- newEx
			//todo add broadcast to user here
		}
	}()

	removeExTicker := time.NewTicker(10 * time.Second)
	go func() {
		for t := range removeExTicker.C {
			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
			var value Exception
			exceptionsMap.RLock()
			for key := range exceptionsMap.m {
				value = exceptionsMap.m[key]
				exceptionsMap.RUnlock()
				break
			}
			value.Show = false
			broadcastException <- value
			//exceptionsMap.Lock() //take the write lock
			//delete(exceptionsMap, value.p)
			//exceptionsMap.Unlock() //take the write lock
			fmt.Println(value)

		}
	}()

	// wait for 10 seconds
	//time.Sleep(20 *time.Second)
	//	ticker.Stop()

}

func main() {
	fmt.Println("websockets project")
	// Configure websocket route

	http.HandleFunc("/ws", handleConnections)
	go broadcastMessages()
	//go exceptiosMapHandler()
	// Start the server on localhost port 8080 and log any errors
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
