package main

import (
	"backendMetaverse/http/client"
	"backendMetaverse/prisma/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

func getRandomeString(lent int) string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, lent)
	for i := range result {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}

type User struct {
	Id      string
	UserId  string
	SpaceId string
	X       float64
	Y       float64
	ws      *websocket.Conn
}

func NewUser() *User {
	u := &User{
		Id: getRandomeString(10),
		X:  0,
		Y:  0,
	}

	return u
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func initHandler(u *User) {
	fmt.Println("In the init handler")

	for {
		_, msg, err := u.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Printf("User %s disconnected: %v", u.Id, err)
			}
			fmt.Println(err)
			u.Close() // Clean up user from rooms, etc.
			break     // Exit the loop to stop handling messages for this user
		}
		var message IncomingMessage
		// fmt.Println(msg)
		json.Unmarshal(msg, &message)
		fmt.Println("The message is ", message)
		switch message.Type {
		case "join":
			u.HandleJoin(message.Payload)
		case "move":
			u.HandleMove(message.Payload)
		}

	}

}

type JoinPayload struct {
	SpaceID string `json:"spaceId"`
	Token   string `json:"token"`
}

type OutgoingMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (u *User) Close() {
	if u.SpaceId != "" {
		GetRoomManager().broadcast(OutgoingMessage{
			Type: "user-left",
			Payload: map[string]interface{}{
				"userId": u.UserId,
			},
		}, u, *&u.SpaceId)
		GetRoomManager().RemoveUser(u, *&u.SpaceId)
	}
	u.ws.Close()
}

func (u *User) HandleJoin(msg json.RawMessage) {
	fmt.Println("In the handle join event handler")
	var joinData JoinPayload
	if err := json.Unmarshal(msg, &joinData); err != nil {
		log.Println(err)
		return
	}
	token, err := jwt.Parse(joinData.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil // Use the same secret as HTTP
	})
	if err != nil || !token.Valid {
		u.Send(OutgoingMessage{Type: "error", Payload: "Invalid token"})
		u.Close()
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	username, ok := claims["username"].(string) // Use the same claim as HTTP
	if !ok {
		u.Send(OutgoingMessage{Type: "error", Payload: "Invalid user ID"})
		u.Close()
		return
	}
	// u.UserId = username
	fmt.Println("REACHED 100")
	prismaClient := client.GetClient()
	fmt.Println(joinData.SpaceID)
	userFound, err := prismaClient.User.FindUnique(db.User.Username.Equals(username)).Exec(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	u.UserId = userFound.ID
	space, err := prismaClient.Space.FindUnique(db.Space.ID.Equals(joinData.SpaceID)).Exec(context.Background())
	fmt.Println("USer id is", u.UserId)
	// fmt.Println(space)
	if err != nil {
		fmt.Println(err)

		u.Send(OutgoingMessage{Type: "error", Payload: "Invalid space"})
		return
	}
	u.SpaceId = joinData.SpaceID
	u.X = rand.Float64() * float64(space.Width)
	u.Y = rand.Float64() * float64(space.Height)
	rm := GetRoomManager()
	rm.addUser(u, u.SpaceId)
	usersInRoom := rm.GetAllUsersInARoom(u.SpaceId)
	filteredUsers := make([]map[string]interface{}, 0)
	allUsers := make([]map[string]interface{}, 0)
	for _, user := range usersInRoom {
		// Include all users (including themselves)
		allUsers = append(allUsers, map[string]interface{}{
			"userId": user.UserId,
			"x":      user.X,
			"y":      user.Y,
		})
	}

	for _, user := range usersInRoom {
		if user.Id != u.Id {
			filteredUsers = append(filteredUsers, map[string]interface{}{"id": user.Id})
		}

	}
	fmt.Println("going to broadcasr")
	u.Send(OutgoingMessage{
		Type: "welcome",
		Payload: map[string]interface{}{
			"userId": u.UserId,
			"x":      u.X,
			"y":      u.Y,
			"users":  allUsers, // Includes all users in the room
		},
	})
	rm.broadcast(OutgoingMessage{
		Type: "user-joined",
		Payload: map[string]interface{}{
			"userId": u.UserId,
			"x":      u.X,
			"y":      u.Y,
		},
	}, u, u.SpaceId)

}

type MovePayload struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (u *User) HandleMove(payload json.RawMessage) {
	var moveData MovePayload
	if err := json.Unmarshal(payload, &moveData); err != nil {
		fmt.Println("error in parsing out json.Rawmessage", err)
		return
	}
	fmt.Println("in the hanldemove ws  func u.x is and moveData.x is", u.X, moveData.X)
	xDisplaceement := math.Abs(u.X - moveData.X)
	yDisplaceement := math.Abs(u.Y - moveData.Y)
	if (xDisplaceement == 1 && yDisplaceement == 0) || (xDisplaceement == 0 && yDisplaceement == 1) {
		u.X = moveData.X
		u.Y = moveData.Y
		fmt.Println("Going to the broadcast func")
		fmt.Println("SPace id is ", u.SpaceId)
		GetRoomManager().broadcast(OutgoingMessage{
			Type: "movement",
			Payload: map[string]interface{}{
				"userId": u.UserId, // <-- Add this
				"x":      u.X,
				"y":      u.Y,
			},
		}, u, u.SpaceId)

	} else {
		u.Send(OutgoingMessage{
			Type: "movement-rejected",
			Payload: map[string]float64{
				"x": u.X,
				"y": u.Y,
			},
		})

	}

}
func (u *User) Send(msg OutgoingMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	if err := u.ws.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
