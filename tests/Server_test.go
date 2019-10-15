package tests

import (
  "gotest.tools/assert"
	"testing"
	"net/http"
	"bytes"
	"time"
	"io"
  srv "server"
)

func postHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK")
}

func Test_Mux_Server_GET(t *testing.T) {
	// Arrange
	var server = srv.NewMuxServer("5001")

	server.MapRoutes([]*srv.Route{
		{"/testget", getHandler, []string{"GET"}}});

	// Act
	go func(){
		server.Start()
	}()

	time.Sleep(10 * time.Second)

	_, err := http.Get("http://127.0.0.1:5001/testget")

	// Assert
	assert.Equal(t, err, nil)
}

func Test_Mux_Server_POST(t *testing.T) {
	// Arrange
	var server = srv.NewMuxServer("5002")

	server.MapRoutes([]*srv.Route{
		{"/testpost", postHandler, []string{"POST"}}});

	// Act
	go func(){
		server.Start()
	}()
	time.Sleep(10 * time.Second)

	_, err := http.Post("http://127.0.0.1:5002/testpost", "", bytes.NewBuffer([]byte{}))

	// Assert
	assert.Equal(t, err, nil)
}
