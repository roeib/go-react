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
	Width  int64 // "0" || "10" ||  "-10"
	Height int64 // "0" || "10" ||  "-10"
}

type point struct {
	X int64  `json:"x,omitempty"`
	Y int64  `json:"y,omitempty"`
}

type Player struct {
	Id            uuid.UUID `json:"id,omitempty"`
	ExceptionType string    `json:"exceptionType,omitempty"`
	Color         [3]int    `json:"color,omitempty"`
	P             *point     `json:"p,omitempty"`
	Show          bool      `json:"show,omitempty"`
	windowH       int64
	windowW       int64
	Collision     bool `json:"collision,omitempty"`
	Score         int  `json:"score,omitempty"`
}

type Exception struct {
	Id            uuid.UUID `json:"id,omitempty"`
	ExceptionType string    `json:"exceptionType,omitempty"`
	Show          bool      `json:"show,omitempty"`
	Color         [3]int    `json:"color,omitempty"`
	P             *point     `json:"p,omitempty"`
}

type ElementsMsg struct {
	Self     *Player    `json:"self,omitempty"`
	Plyer    *Player    `json:"player,omitempty"`
	Excption *Exception `json:"exception,omitempty"`
}

var exceptionsMap = struct {
	sync.RWMutex
	m map[point]Exception
}{m: make(map[point]Exception)}

var s = rand.NewSource(time.Now().UnixNano())
var exceptionsTypes = [3]string{"NullPointerException", "DivideByZeroException", "IOException"}
var clients = make(map[*websocket.Conn]*Player) // connected clients
var broadcastMsg = make(chan ElementsMsg)
var upgrader = websocket.Upgrader{}

func handleNewPlayer(ws *websocket.Conn) {
	rand.Seed(time.Now().UnixNano())
	player := Player{Id: uuid.New(), P: &point{X: int64(rand.Intn(300)), Y: int64(rand.Intn(300))}, Score: 0, Show: true, ExceptionType: exceptionsTypes[rand.Intn(3)], Color: [3]int{rand.Intn(256), rand.Intn(256), rand.Intn(256)}, Collision: false}
	fmt.Println("new player")
	fmt.Println(player)

	//send to client active player as self
	m := ElementsMsg{Self: &player, Plyer: &player}
	ws.WriteJSON(m)
	//send to new player all current players
	for key := range clients {
		m := ElementsMsg{Plyer: clients[key]}
		err := ws.WriteJSON(m)
		if err != nil {
			log.Printf("84 error: %v", err)
			ws.Close()
			delete(clients, ws)
		}
	}

	//send to new player all current exceptions
	for key := range exceptionsMap.m {
		var e = exceptionsMap.m[key]
		m := ElementsMsg{Excption: &e}
		err := ws.WriteJSON(m)
		if err != nil {
			log.Printf("96 error: %v", err)
			ws.Close()
			delete(clients, ws)
		}
	}

	//broadcast new player to all clients
	ms := ElementsMsg{Plyer: &player}
	broadcastMsg <- ms

	//update player window Width/Height
	clients[ws] = &player
	var msg screenWH
	err := ws.ReadJSON(&msg)
	if err != nil {
		log.Printf(" 89 error: %v", err)
		var plyr = clients[ws]
		plyr.Show = false
		ms := ElementsMsg{Plyer: plyr}
		broadcastMsg <- ms
		delete(clients, ws)
	}
	clients[ws].windowH = msg.Height
	clients[ws].windowW = msg.Width
}

func handlePlayerMovement(ws *websocket.Conn, newX int64, newY int64) {
	x := int64(clients[ws].P.X) + newX
	y := int64(clients[ws].P.Y) + newY
	player := *clients[ws]
	if y < 0 || x < 0 || x >= clients[ws].windowW || y >= clients[ws].windowH {
		player.Collision = true
	} else {
		for key := range exceptionsMap.m {
			value := exceptionsMap.m[key]
			if (value.ExceptionType == player.ExceptionType) &&
			   (x == value.P.X || x+50 >= value.P.X || x-50 <= value.P.X) &&
				(y == value.P.Y || y+50 >= value.P.Y || y-50 <= value.P.Y) {
				fmt.Println("Ex found is: ", value)
				value.Show = false
				ms := ElementsMsg{Excption: &value}
				broadcastMsg <- ms
				delete(exceptionsMap.m, *value.P)
				player.Score = player.Score + 1
				clients[ws].Score = player.Score
				break
			}
		}
		player.P.X = x
		player.P.Y = y
		clients[ws].P.X = x
		clients[ws].P.Y = y
	}
	ms := ElementsMsg{Plyer: &player}
	broadcastMsg <- ms

}

func exceptionMapHandler(){
	time.Sleep(30 *time.Second)
	rand.Seed(time.Now().UnixNano())
	min := 50
	max := 300
	var newEx = Exception{Id: uuid.New(), ExceptionType: exceptionsTypes[rand.Intn(3)], P: &point{X: int64(rand.Intn(max - min + 1) + min), Y: int64(rand.Intn(max - min + 1) + min)}, Show: true, Color: [3]int{0, 0, 0}}
	exceptionsMap.m[*newEx.P] = newEx
	ms := ElementsMsg{Excption: &newEx}
	broadcastMsg <- ms
	fmt.Println("added EX element")
	fmt.Println(newEx)
	time.Sleep(50 *time.Second)
	newEx.Show = false
	ms = ElementsMsg{Excption: &newEx}
	broadcastMsg <- ms
	delete(exceptionsMap.m, *newEx.P)
}

func exceptionsMapHandler() {
	time.Sleep(30 *time.Second)
	rand.Seed(time.Now().UnixNano())
	min := 50
	max := 300
	addExTicker := time.NewTicker(40 * time.Second)
	go func() {
		for t := range addExTicker.C {
			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
			var newEx = Exception{Id: uuid.New(), ExceptionType: exceptionsTypes[rand.Intn(3)], P: &point{X: int64(rand.Intn(max - min + 1) + min), Y: int64(rand.Intn(max - min + 1) + min)}, Show: true, Color: [3]int{0, 0, 0}}
			exceptionsMap.m[*newEx.P] = newEx
			ms := ElementsMsg{Excption: &newEx}
			broadcastMsg <- ms
			fmt.Println("added EX element")
			fmt.Println(newEx)
		}
	}()

	removeExTicker := time.NewTicker(50 * time.Second)
	go func() {
		for t := range removeExTicker.C {
			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
			var value Exception
			for key := range exceptionsMap.m {
				value = exceptionsMap.m[key]
				break
			}
			value.Show = false
			ms := ElementsMsg{Excption: &value}
			broadcastMsg <- ms
			value, ok := exceptionsMap.m[*value.P]
			if ok {
				delete(exceptionsMap.m, *value.P)
			}
			fmt.Println("removed EX element")
			fmt.Println(value)
		}
	}()

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
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("206 error: %v", err)
			plyrMsg := clients[ws]
			plyrMsg.Show = false
			ms := ElementsMsg{Plyer: plyrMsg}
			broadcastMsg <- ms
			delete(clients, ws)
			break
		}
		newX, _ := strconv.ParseInt(msg.X, 10, 64)
		newY, _ := strconv.ParseInt(msg.Y, 10, 64)
		handlePlayerMovement(ws, newX, newY)
	}
}
func broadcastMessages() {
	for {
		msg := <-broadcastMsg
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("226 error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	fmt.Println("ExceptionalMonkeys ... ")
	http.HandleFunc("/ws", handleConnections)
	go broadcastMessages()
	 go exceptionsMapHandler()
  //  go exceptionMapHandler()
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
