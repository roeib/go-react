package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Define our message object
type actionMessage struct {
	x string
	y string
}

type point struct{
	x        int64
	y        int64
}

type Player struct {
	exceptionType string
	color [3]int
	size int
	p point
	show bool

}

type Exception struct {
	exceptionType string
	show bool
	p point
}

var exceptionsMap = struct{
	sync.RWMutex
	m map[point]Exception
}{m: make(map[point]Exception)}

var s = rand.NewSource(time.Now().UnixNano())
var exceptionsTypes = [3]string{"NullPointerException", "DivideByZeroException", "IOException"}
var clients = make(map[*websocket.Conn]*Player) // connected clients
var broadcastPlayers = make(chan Player)
var broadcastException = make(chan Exception)
var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	 r2 := rand.New(s)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //no cors

	ws, err := upgrader.Upgrade(w, r, nil) // Upgrade initial GET request to a websocket
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close() // Make sure we close the connection when the function returns

	player := Player{ p:point{x:0, y:0}, size: 50, show:true, exceptionType: exceptionsTypes[rand.Intn(3)], color:[3]int{r2.Intn(256), r2.Intn(256), r2.Intn(256)} }
	clients[ws] = &player

	fmt.Println("new player")
	fmt.Println(clients[ws])
	for {
		var msg actionMessage
		// Read in a new message as JSON and map it to a incomingMessage object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			var plyrMsg =  clients[ws]
			plyrMsg.show= false;
			broadcastPlayers <- *plyrMsg
			delete(clients, ws)
			break
		}
		fmt.Println("received msg:")
		fmt.Println(msg)

		fmt.Println("player curr values")
		fmt.Println(clients[ws])
		//INSERT HERE  call to check collision function
		newX, _ :=  strconv.ParseInt(msg.x, 10, 64)
		clients[ws].p.x = int64(clients[ws].p.x) + newX

	//	newY, _ :=  strconv.ParseInt(msg.y, 10, 64)
	//	clients[ws].p.y = int64(clients[ws].p.y) + newY
		fmt.Println("player with new values ")
		fmt.Println(clients[ws])
		broadcastPlayers <- *clients[ws]
	}
}
func broadcastMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcastPlayers

		fmt.Println("!!!!!broadcast msg:")
		fmt.Println(msg)

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



func exceptiosMapHandler (){
	//Thread.
	// every 5 sec:  create new ex ->add to eXarr -> broadcast to users
	//exery 10 sec:  choose rand ex ->remove from exArr ->broadcast to users

	var r = rand.New(s)
	addExTicker := time.NewTicker(5* time.Second)
	go func() {
		for t := range addExTicker.C {
			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
			var newEx = Exception{exceptionType: exceptionsTypes[r.Intn(3)], p:point{x:int64(r.Intn(255)), y:int64(r.Intn(255))}, show:true }
			exceptionsMap.Lock() //take the write lock
			exceptionsMap.m[newEx.p]=newEx;
			exceptionsMap.Unlock() // release the write lock
			fmt.Println(newEx)
			broadcastException <- newEx
			//todo add broadcast to user here
		}
	}()

	removeExTicker := time.NewTicker(10* time.Second)
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
			value.show =false
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
