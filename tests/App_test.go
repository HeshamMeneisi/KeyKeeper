package tests

import (
	kk "app"
	"bytes"
	cfg "config"
	db "dbmgr"
	"encoding/json"
	"gotest.tools/assert"
	"net/http"
	s "server"
	"testing"
	"time"
)

const (
	PORT    = "5000"
	KEY_TTL = time.Hour
)

func Test_Simulation(t *testing.T) {
	// Arrange

	// Load config
	config := cfg.GetConfig("../config.yml")
	// Connect to DB
	client, _ := db.NewMongoClient(config.Database.Host, config.Database.Port)
	// Start a server
	server := s.NewMuxServer(PORT)
	// Create the app
	app := kk.NewKeyKeeper("kkDB", client)
	app.RegisterRoutes(server)

	// Act
	go func() {
		server.Start()
	}()

	time.Sleep(10 * time.Second)

	message := map[string]string{
		"key": "TEST_KEY"}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		t.Errorf("%v", err)
	}

	resp, err := http.Post("http://127.0.0.1:5000/keys", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		t.Errorf("%v", err)
	}
	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	//   if err != nil {
	//       t.Log(err)
	//   }
	//   bodyString := string(bodyBytes)
	// t.Log(bodyString)
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	var id = result["id"]
	t.Log("Received ID:", id)

	resp, err = http.Get("http://127.0.0.1:5000/key/" + string(id))
	if err != nil {
		t.Errorf("%v", err)
	}

	json.NewDecoder(resp.Body).Decode(&result)

	t.Log("Received Key:", result["key"])
	assert.Equal(t, "TEST_KEY", result["key"])
}
