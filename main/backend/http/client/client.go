package client

import (
	"backendMetaverse/prisma/db"
	"sync"
)

// singleton variables
var (
	clientInstance *db.PrismaClient
	once           sync.Once
)

func GetClient() *db.PrismaClient {
	once.Do(func() {
		clientInstance = db.NewClient()
	})
	return clientInstance
}
