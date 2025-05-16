package main

import (
	"backendMetaverse/http/client"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgradation failed", err)
		return
	}
	fmt.Println("A client connected")
	user := NewUser()
	// fmt.Println(user)
	user.ws = conn       // Critical missing piece!
	go initHandler(user) // Start message processing
}

func main() {
	fmt.Println("Trying to connect to the ws")
	r := mux.NewRouter()
	PrismaClient := client.GetClient()

	if err := PrismaClient.Prisma.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Connected successfullt")
	defer PrismaClient.Prisma.Disconnect()
	r.HandleFunc("/", wsHandler)
	err := http.ListenAndServe(":3002", r)
	defer fmt.Println("Connected Successfully`")
	if err != nil {
		fmt.Println(err)
	}

}
