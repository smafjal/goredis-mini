package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
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
			response := executePING()
			conn.Write([]byte(response))
		case "SET":
			response := executeSET(cmd)
			conn.Write([]byte(response))
		case "SETEX":
			response := executeSETEX(cmd)
			conn.Write([]byte(response))
		case "EXPIRE":
			response := executeEXPIRE(cmd)
			conn.Write([]byte(response))
		case "GET":
			response := executeGET(cmd)
			conn.Write([]byte(response))
		case "DEL":
			response := executeDEL(cmd)
			conn.Write([]byte(response))
		default:
			conn.Write([]byte("-unknown command\r\n"))
		}
	}
}

func executePING() string {
	return "+PONG\r\n"
}

func executeSET(cmd []string) string {
	if len(cmd) != 3 {
		return "-wrong number of arguments for `set`\r\n"
	}
	store.Set(cmd[1], cmd[2])
	return "+OK\r\n"
}

func executeSETEX(cmd []string) string {
	if len(cmd) != 4 {
		return "-wrong number of arguments for `setex`\r\n"
	}
	ttl, err := strconv.Atoi(cmd[2])
	if err != nil {
		return "invalid ttl\r\n"
	}
	store.SetWithTTL(cmd[1], cmd[3], ttl)
	return "+OK\r\n"
}

func executeGET(cmd []string) string {
	if len(cmd) != 2 {
		return "-wrong number of arguments for `get`\r\n"
	}
	if value, ok := store.Get(cmd[1]); ok {
		return fmt.Sprintf("$ %s\r\n", value)
	}
	return "$ -1\r\n"
}

func executeDEL(cmd []string) string {
	if len(cmd) != 2 {
		return "-wrong number of arguments for `del`\r\n"
	}
	ok := store.Del(cmd[1])
	return fmt.Sprintf(": %d\r\n", ok)
}

func executeEXPIRE(cmd []string) string {
	if len(cmd) != 3 {
		return "-wrong number of arguments for `expire`\r\n"
	}
	ttl, err := strconv.Atoi(cmd[2])
	if err != nil {
		return "invalid ttl\r\n"
	}

	if value, ok := store.Get(cmd[1]); ok {
		store.SetWithTTL(cmd[1], value, ttl)
		return ": 1\r\n"
	}
	return ": 0\r\n"
}
