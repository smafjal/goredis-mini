package main

import (
	"fmt"

	"github.com/smafjal/goredis-mini/internal/core"
	"github.com/smafjal/goredis-mini/internal/server"
	"github.com/smafjal/goredis-mini/internal/store"
)

var address = ":8980"

func main() {
	fmt.Println("goredis-mini starting at", address)

	aoflogger := store.NewAof("aof/appendonly.aof")
	db := store.NewDatabase(aoflogger)
	eng := core.NewEngine(db)

	lines, err := aoflogger.Load()
	if err == nil {
		for _, line := range lines {
			eng.ProcessAofCmd(line, db)
		}
	}
	server := server.NewServer(address, eng)
	server.Start()
}
