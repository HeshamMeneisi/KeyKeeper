package keykeeper

import (
	"context"
	db "dbmgr"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
	"os"
	srv "server"
	"sync/atomic"
	"time"
)

type KKApp struct {
	keyCol   *mongo.Collection
	dbClient *mongo.Client
	cid      CustomID
	healthy  int32
	stats    Stats
	ttl      time.Duration
}

type Stats struct {
	KSinceStarted uint32 `json:"keys_since_started"`
	KLastHour     uint32 `json:"keys_last_hour"`
}

type Key struct {
	Value     string    `json:"key" bson:"value"`
	Timestamp time.Time `bson:"timestamp"`
	ID        uint32    `bson:"_id"`
}

func NewKeyKeeper(dbName string, dbClient *mongo.Client) *KKApp {
	app := new(KKApp)
	app.dbClient = dbClient
	app.keyCol = dbClient.Database(dbName).Collection("keys")
	app.stats = Stats{0, 0}
	c := db.NextIdForCol(app.keyCol)
	log.Println("Starting Index", c)
	app.cid = NewBase64ID(c)
	atomic.StoreInt32(&app.healthy, 1)
	// Health monitoring
	quit := make(chan os.Signal, 1)
	go func() {
		<-quit
		log.Println("Server is shutting down...")
		atomic.StoreInt32(&app.healthy, 0)

		// Attempt recovery?
		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}()
	return app
}

func (app *KKApp) ScheduleCleanup(d time.Duration) {
	app.ttl = d
	// Cleanup every d
	// No need for CRON jobs if we cleanup on server restart
	app.cleanupOldKeys()
	go func() {
		for _ = range time.Tick(d) {
			app.cleanupOldKeys()
		}
	}()
}

func (app *KKApp) RegisterRoutes(server srv.Server) {
	server.MapRoutes([]*srv.Route{
		{"/keys", app.createHandler(addKey), []string{"POST"}},
		{"/key/{id}", app.createHandler(getKey), []string{"GET"}},
		{"/health_check", app.createHandler(checkHealth), []string{"GET"}}})
}

func (app *KKApp) createHandler(f func(*KKApp, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(app, w, r)
	}
}

func (app *KKApp) cleanupOldKeys() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"timestamp": bson.M{
		"$lt": time.Now().Add(-app.ttl),
	}}
	_, err := app.keyCol.DeleteMany(ctx, filter)
	if err != nil {
		log.Println("Failed to cleanup old keys.")
		log.Println(err)
	}
	app.stats.KLastHour = 0
}

func addKey(app *KKApp, w http.ResponseWriter, r *http.Request) {
	// Retrieve request content
	w.Header().Set("Content-Type", "application/json")
	var obj Key
	if r.Body == nil {
		http.Error(w, "JSON object missing.", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Get an ID for the key
	obj.Timestamp = time.Now()
	obj.ID = app.cid.GetNextUint32()

	// Insert into DB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = app.keyCol.InsertOne(ctx, obj)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	app.stats.KSinceStarted += 1
	app.stats.KLastHour += 1

	// Formulate response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		struct {
			Exported string `json:"id"`
		}{app.cid.FromUint32(obj.ID)})
}

func getKey(app *KKApp, w http.ResponseWriter, r *http.Request) {
	// Retrieve request parameters
	params := mux.Vars(r)
	b64 := params["id"]
	oid, err := app.cid.ToUint32(b64)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Query DB
	var obj Key
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"_id": oid}
	err = app.keyCol.FindOne(ctx, filter).Decode(&obj)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Formulate response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Exported string `json:"key"`
	}{obj.Value})
}

func checkHealth(app *KKApp, w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&app.healthy) == 1 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(app.stats)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, "Service DOWN")
	}
}
