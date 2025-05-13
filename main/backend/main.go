package main

import (
	"backendMetaverse/client"
	"backendMetaverse/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("In the main file")

	PrismaClient := client.GetClient()
	// Initialize Prisma client
	if err := PrismaClient.Prisma.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer PrismaClient.Prisma.Disconnect()

	// Optionally: Pass client to router if needed
	// r := router.NewRouter(client)
	// log.Fatal(http.ListenAndServe(":3000", r))

	// Start HTTP server
	log.Println("Server started on :3000")
	if err := http.ListenAndServe(":3000", router.Router()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
