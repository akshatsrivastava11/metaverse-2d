package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type RoomManager struct {
	rooms map[string][]*User
}

var (
	instance *RoomManager
	once     sync.Once
)

func GetRoomManager() *RoomManager {
	once.Do(func() {
		instance = &RoomManager{
			rooms: make(map[string][]*User),
		}
	})
	return instance
}
func (rm *RoomManager) addUser(user *User, spaceId string) {
	users := rm.rooms[spaceId]
	users = append(users, user)
	rm.rooms[spaceId] = users
}

func (rm *RoomManager) RemoveUser(userTobeDeleted *User, spaceId string) {
	users := rm.rooms[spaceId]
	for ind, user := range users {
		if user.Id == userTobeDeleted.Id {
			removed := users[0:ind]
			removed = append(removed, users[ind+1:]...)
			rm.rooms[spaceId] = removed
			break
		}
	}
}

func (rm *RoomManager) GetAllUsersInARoom(spaceId string) []*User {
	users := rm.rooms[spaceId]
	return users
}
func (rm *RoomManager) broadcast(message OutgoingMessage, userThatSendsTheMessage *User, spaceId string) {
	fmt.Println("In the broadcst fun")
	users, exist := rm.rooms[spaceId]
	if !exist {
		fmt.Println("No such space exist bruhhh")
	}
	val, _ := json.Marshal(message)
	for _, user := range users {
		// Remove this check: if user.Id != userThatSendsTheMessage.Id
		user.ws.WriteMessage(websocket.TextMessage, val)
	}
}
