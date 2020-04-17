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

type Player struct {
	Id            uuid.UUID `json:"id"`
	ExceptionType string    `json:"exceptionType"`
	Color         [3]int    `json:"color"`
	X int64  `json:"x"`
	Y int64  `json:"y"`
	Show          bool      `json:"show"`
	windowH       int64
	windowW       int64
	Collision     bool `json:"collision"`
	Score         int  `json:"score"`
}

type Exception struct {
	Id            uuid.UUID `json:"id"`
	ExceptionType string    `json:"exceptionType"`
	Show          bool      `json:"show"`
	Color         [3]int    `json:"color"`
	X int64  `json:"x"`
	Y int64  `json:"y"`
}

type ElementsMsg struct {
	Self     *Player    `json:"self,omitempty"`
	Plyer    *Player    `json:"player,omitempty"`
	Excption *Exception `json:"exception,omitempty"`
}

var exceptionsMap = struct {
	sync.RWMutex
	items [3]Exception
}{}

//todo amount of exceptions according to amount of possible exceptions?
func initExceptionsList(){
	exceptionsMap.Lock()
	defer exceptionsMap.Unlock()

	for j := 0; j < 3; j++ {
		exceptionsMap.items[j] = Exception{Id: uuid.New(), ExceptionType: exceptionsTypes[rand.Intn(3)], X: 0, Y: 0, Show: false}
	}
}

func Set()Exception {
	rand.Seed(time.Now().UnixNano())
	min := 50
	max := 300
	exceptionsMap.Lock()
	defer exceptionsMap.Unlock()

	var indx =0

	for j := 0; j < 3; j++ {
		if !(exceptionsMap.items[j].Show){
			exceptionsMap.items[j].Show = true
			exceptionsMap.items[j].X= int64(rand.Intn(max - min + 1) + min)
			exceptionsMap.items[j].Y= int64(rand.Intn(max - min + 1) + min)
			indx = j
			break
		}
	}
	//	fmt.Println("timer: added EX element")
	//fmt.Println(value)
	var res = exceptionsMap.items[indx]
	return res
}

func RemoveRand() Exception {
	exceptionsMap.Lock()
	defer exceptionsMap.Unlock()

	var indx =0

	for j := 0; j < 3; j++ {
		if !(exceptionsMap.items[j].Show){
			exceptionsMap.items[j].Show = false
			indx = j
			break
		}
	}
	//	fmt.Println("timer: added EX element")
	//fmt.Println(value)
	var res = exceptionsMap.items[indx]
	return res
}

func RemoveIfPossible(newX int64, newY int64 ,player Player ) (Exception, bool) {
	exceptionsMap.Lock()
	defer exceptionsMap.Unlock()
	var value Exception
	for j:= 0; j< 3; j++{
		value=exceptionsMap.items[j]
		if value.Show {
			if (value.ExceptionType == player.ExceptionType) &&
			(newX == value.X || newX+50 >= value.X || newX-50 <= value.X) &&
			(newY == value.Y || newY+50 >= value.Y || newY-50 <= value.Y) {
			fmt.Println("Ex found is: ", value)
			value.Show = false
			exceptionsMap.items[j].Show=false
			return value, true
			}
		}
	}
	return value, false
}
var s = rand.NewSource(time.Now().UnixNano())
var exceptionsTypes = [3]string{"NullPointerException", "DivideByZeroException", "IOException"}
var clients = make(map[*websocket.Conn]*Player) // connected clients
var broadcastMsg = make(chan ElementsMsg)
var upgrader = websocket.Upgrader{}

func handleNewPlayer(ws *websocket.Conn) {
	rand.Seed(time.Now().UnixNano())
	player := Player{Id: uuid.New(), X: int64(rand.Intn(300)), Y: int64(rand.Intn(300)), Score: 0, Show: true, ExceptionType: exceptionsTypes[rand.Intn(3)], Color: [3]int{rand.Intn(256), rand.Intn(256), rand.Intn(256)}, Collision: false}
	fmt.Println("new player")
	fmt.Println(player)

	//send to client active player as self
	m := ElementsMsg{Self: &player}
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
	x := int64(clients[ws].X) + newX
	y := int64(clients[ws].Y) + newY
	player := *clients[ws]
	if y < 0 || x < 0 || x >= clients[ws].windowW || y >= clients[ws].windowH {
		player.Collision = true
	} else {
		value, ok:= RemoveIfPossible(x,y,player)
		if ok {
			ms := ElementsMsg{Excption: &value}
			broadcastMsg <- ms
			player.Score = player.Score + 1
			clients[ws].Score = player.Score
		}
		player.X = x
		player.Y = y
		clients[ws].X = x
		clients[ws].Y = y
	}
	ms := ElementsMsg{Plyer: &player}
	broadcastMsg <- ms
	fmt.Println("msg sent with new score: ")
	fmt.Println(ms)
}

func exceptionsMapHandler() {
	time.Sleep(30 *time.Second)

	addExTicker := time.NewTicker(40 * time.Second)
	go func() {
		for t := range addExTicker.C {
			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
			newEx := Set()
			ms := ElementsMsg{Excption: &newEx}
			broadcastMsg <- ms
		}
	}()

	removeExTicker := time.NewTicker(50 * time.Second)
	go func() {
		for t := range removeExTicker.C {
			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
			value :=RemoveRand()
			ms := ElementsMsg{Excption: &value}
			broadcastMsg <- ms
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
	initExceptionsList()
	go broadcastMessages()
	go exceptionsMapHandler()

	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
