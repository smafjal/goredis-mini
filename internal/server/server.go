package server

import (
	"bufio"
	"net"
	"strings"

	"github.com/smafjal/goredis-mini/internal/core"
)

type Server struct {
	address string
	eng     *core.Engine
}

func NewServer(address string, eng *core.Engine) *Server {
	return &Server{
		address: address,
		eng:     eng,
	}
}

func (s *Server) Start() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, s.eng)
	}
}

func handleConnection(conn net.Conn, eng *core.Engine) {
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
			response := eng.ExecutePING()
			conn.Write([]byte(response))
		case "SET":
			response := eng.ExecuteSET(cmd)
			conn.Write([]byte(response))
		case "SETEX":
			response := eng.ExecuteSETEX(cmd)
			conn.Write([]byte(response))
		case "EXPIRE":
			response := eng.ExecuteEXPIRE(cmd)
			conn.Write([]byte(response))
		case "GET":
			response := eng.ExecuteGET(cmd)
			conn.Write([]byte(response))
		case "DEL":
			response := eng.ExecuteDEL(cmd)
			conn.Write([]byte(response))
		default:
			conn.Write([]byte("$ unknown command\r\n"))
		}
	}
}
