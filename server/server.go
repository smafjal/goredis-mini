package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/smafjal/goredis-mini/db"
)

var store = db.NewDatabase()

func Start(address string) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		cmd := strings.Fields(strings.TrimSpace(line))
		if len(cmd) == 0 {
			continue
		}

		switch strings.ToUpper(cmd[0]) {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "SET":
			if len(cmd) != 3 {
				conn.Write([]byte("-wrong number of arguments for `set`\r\n"))
			}
			store.Set(cmd[1], cmd[2])
			conn.Write([]byte("+OK\r\n"))
		case "GET":
			if len(cmd) != 2 {
				conn.Write([]byte("-wrong number of arguments for `get`\r\n"))
			} else {
				if value, ok := store.Get(cmd[1]); ok {
					msg := fmt.Sprintf("$ %s\r\n", value)
					conn.Write([]byte(msg))
				} else {
					conn.Write([]byte("$ -1\r\n"))
				}
			}
		case "DEL":
			if len(cmd) != 2 {
				conn.Write([]byte("-wrong number of arguments for `del`\r\n"))
			} else {
				ok := store.Del(cmd[1])
				conn.Write([]byte(fmt.Sprintf(": %d\r\n", ok)))
			}
		default:
			conn.Write([]byte("-unknown command\r\n"))
		}
	}
}
