package tests

import (
	"testing"
	"gotest.tools/assert"
	kk "app"
	db "dbmgr"
	s "server"
	"net/http"
	"bytes"
	"time"
	"encoding/json"
)

const (
	PORT    = "8000"
	KEY_TTL = time.Hour
)

func Test_Simulation(t *testing.T) {
	// Arrange
	client, _ := db.NewMongoClient("") // Docker: mongodb
	server := s.NewMuxServer(PORT)
	app := kk.NewKeyKeeper("kkDB", client)
	app.RegisterRoutes(server)

	// Act
	server.Start()

	message := map[string]string{
		"key": "TEST_KEY"	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		t.Errorf("%v", err)
	}

	resp, err := http.Post("http://127.0.0.1:8000/keys", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Errorf("%v", err)
	}
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	var id = result["id"]

	resp, err = http.Get("http://127.0.0.1:8000/key/"+string(id))
	if err != nil {
		t.Errorf("%v", err)
	}

	json.NewDecoder(resp.Body).Decode(&result)

	assert.Equal(t, "TEST_KEY", result["key"])
}
