package main

import (
	"backendMetaverse/http/client"
	"backendMetaverse/http/router"
	"fmt"
	"log"
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
func main() {
	fmt.Println("In the main file")

	PrismaClient := client.GetClient()
	handler := CORS(router.Router())
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
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
