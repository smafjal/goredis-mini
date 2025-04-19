package benchmark

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
)

var addr = "localhost:8980"

func setupConn(tb testing.TB) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		tb.Fatalf("failed to connect: %v", err)
	}
	return conn
}

func BenchmarkSet(b *testing.B) {
	conn := setupConn(b)
	defer conn.Close()
	reader := bufio.NewReader(conn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench-key-%d", i)
		val := fmt.Sprintf("bench-val-%d", i)
		fmt.Fprintf(conn, "SET %s %s\n", key, val)
		resp, err := reader.ReadString('\n')
		if err != nil || !strings.HasPrefix(resp, "$ OK") {
			b.Fatalf("SET failed: %v, response: %s", err, resp)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	conn := setupConn(b)
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Prepare data first
	fmt.Fprintf(conn, "SET bench-get-key 42\n")
	_, _ = reader.ReadString('\n')

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(conn, "GET bench-get-key\n")
		resp, err := reader.ReadString('\n')
		if err != nil || !strings.Contains(resp, "$ 42") {
			b.Fatalf("GET failed: %v, response: %s", err, resp)
		}
	}
}
