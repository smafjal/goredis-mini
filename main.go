package main

import (
	"fmt"

	"github.com/smafjal/goredis-mini/server"
)

var address = ":8980"

func main() {
	fmt.Println("goredis-mini fun starting at", address)
	server.Start(address)
}
