# goredis-mini

A minimal Redis clone written in Go — lightweight, fast, and built for learning!  
Supports core Redis features like `SET`, `GET`, `DEL`, key expiration, AOF persistence, Pub/Sub, and more.

---

## ✨ Features

- Basic Commands: `PING`, `SET`, `GET`, `DEL`, `EXPIRE`
- Key Expiration: Automatically delete keys after a TTL
- AOF Persistence: Write every command to disk for durability
- Pub/Sub: Publish messages across channels
- 💻 Simple CLI Client: Connects to the server like the real Redis CLI

---

## 🧪 Demo

```bash
$ go run server/main.go
Server started on :8980

$ go run client/cli.go
goredis-mini> PING
+PONG
goredis-mini> SET hello world
+OK
goredis-mini> GET hello
$ world

```

### Happy Learning 📚✨