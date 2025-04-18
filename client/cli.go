package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8980")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		conn.Close()
		os.Exit(0)
	}()

	reader := bufio.NewReader(conn)
	input := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("goredis-mini> ")
		line, _ := input.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		_, err := conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}

		for {
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Disconnected.")
				return
			}
			fmt.Print(response)
			if strings.HasPrefix(response, "$") {
				break
			}
		}
	}
}
