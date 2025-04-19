package core

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/smafjal/goredis-mini/internal/pubsub"
	"github.com/smafjal/goredis-mini/internal/store"
)

type Engine struct {
	DB *store.Database
	PS *pubsub.Pubsub
}

func NewEngine(db *store.Database, ps *pubsub.Pubsub) *Engine {
	return &Engine{
		DB: db,
		PS: ps,
	}
}

func (e *Engine) ProcessAofCmd(line string, db *store.Database) {
	cmd := strings.Fields(strings.TrimSpace(line))
	if len(cmd) == 0 {
		return
	}

	switch strings.ToUpper(cmd[0]) {
	case "SET":
		if len(cmd) == 3 {
			e.ExecuteSET(cmd)
		}
	case "GET":
		if len(cmd) == 2 {
			e.ExecuteGET(cmd)
		}
	case "SETEX":
		if len(cmd) == 4 {
			e.ExecuteSETEX(cmd)
		}
	case "DEL":
		if len(cmd) == 2 {
			e.ExecuteDEL(cmd)
		}
	case "EXPIRE":
		if len(cmd) == 3 {
			e.ExecuteEXPIRE(cmd)
		}
	}
}

func (e *Engine) ExecutePING() string {
	return "$ PONG\r\n"
}

func (e *Engine) ExecuteSET(cmd []string) string {
	if len(cmd) != 3 {
		return "$ wrong number of arguments for `set`\r\n"
	}
	e.DB.Set(cmd[1], cmd[2])
	return "$ OK\r\n"
}

func (e *Engine) ExecuteSETEX(cmd []string) string {
	if len(cmd) != 4 {
		return "$ wrong number of arguments for `setex`\r\n"
	}
	ttl, err := strconv.Atoi(cmd[2])
	if err != nil {
		return "$ invalid ttl\r\n"
	}
	e.DB.SetWithTTL(cmd[1], cmd[3], ttl)
	return "$ OK\r\n"
}

func (e *Engine) ExecuteGET(cmd []string) string {
	if len(cmd) != 2 {
		return "$ wrong number of arguments for `get`\r\n"
	}
	if value, ok := e.DB.Get(cmd[1]); ok {
		return fmt.Sprintf("$ %s\r\n", value)
	}
	return "$ -1\r\n"
}

func (e *Engine) ExecuteDEL(cmd []string) string {
	if len(cmd) != 2 {
		return "$ wrong number of arguments for `del`\r\n"
	}
	ok := e.DB.Del(cmd[1])
	return fmt.Sprintf(": %d\r\n", ok)
}

func (e *Engine) ExecuteEXPIRE(cmd []string) string {
	if len(cmd) != 3 {
		return "$ wrong number of arguments for `expire`\r\n"
	}
	ttl, err := strconv.Atoi(cmd[2])
	if err != nil {
		return "$ invalid ttl\r\n"
	}

	if value, ok := e.DB.Get(cmd[1]); ok {
		e.DB.SetWithTTL(cmd[1], value, ttl)
		return "$ 1\r\n"
	}
	return "$ 0\r\n"
}

func (e *Engine) ExecuteSUBSCRIBE(cmd []string, conn net.Conn) string {
	if len(cmd) < 2 {
		return "$ wrong number of arguments for `subscribe`\r\n"
	}
	e.PS.Subscribe(cmd[1:], conn)
	return "$ OK\r\n"
}

func (e *Engine) ExecutePUBLISH(cmd []string) string {
	if len(cmd) < 3 {
		return "$ wrong number of arguments for `publish`\r\n"
	}
	count := e.PS.Publish(cmd[1], strings.Join(cmd[2:], " "))
	return fmt.Sprintf("$ %d\r\n", count)
}

func (e *Engine) ExecuteUNSUBSCRIBE(cmd []string, conn net.Conn) string {
	if len(cmd) < 2 {
		return "$ wrong number of arguments for `unsubscribe`\r\n"
	}
	e.PS.Unsubscribe(cmd, conn)
	return "$ OK\r\n"
}
