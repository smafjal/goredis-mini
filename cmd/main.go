package main

import (
	"log"
	"net/http"

	"github.com/smafjal/goredis-mini/api"
	"github.com/smafjal/goredis-mini/internal/core"
	"github.com/smafjal/goredis-mini/internal/pubsub"
	"github.com/smafjal/goredis-mini/internal/store"
)

func main() {
	aoflogger := store.NewAof("aof/appendonly.aof")
	db := store.NewDatabase(aoflogger)
	ps := pubsub.NewPubsub()
	eng := core.NewEngine(db, ps)

	handler := api.NewHandler(eng)

	http.HandleFunc("/get", handler.Get)
	http.HandleFunc("/set", handler.Set)

	log.Println("Starting HTTP server on :8980")
	log.Fatal(http.ListenAndServe(":8980", nil))
}
