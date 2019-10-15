package main

import (
	kk "app"
	db "dbmgr"
	s "server"
	"time"
)

const (
	PORT    = "8000"
	KEY_TTL = time.Hour
)

func main() {
	// Connect to DB
	client, _ := db.NewMongoClient("") // Docker: mongodb

	// Setup Server

	server := s.NewMuxServer(PORT)

	// Setup App
	app := kk.NewKeyKeeper("kkDB", client)

	app.RegisterRoutes(server)

	app.ScheduleCleanup(KEY_TTL)

	// Start server
	server.Start()
}
