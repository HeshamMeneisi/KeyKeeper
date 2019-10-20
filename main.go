package main

import (
	kk "app"
	cfg "config"
	db "dbmgr"
	s "server"
	"time"
)

const (
	KEY_TTL = time.Hour
)

func main() {
	// Load config
	config := cfg.GetConfig("config.yml")

	// Connect to DB
	client, _ := db.NewMongoClient(config.Database.Host, config.Database.Port)

	// Setup Server

	server := s.NewMuxServer(config.Server.Port)

	// Setup App
	app := kk.NewKeyKeeper("kkDB", client)

	app.RegisterRoutes(server)

	app.ScheduleCleanup(KEY_TTL)

	// Start server
	server.Start()
}
