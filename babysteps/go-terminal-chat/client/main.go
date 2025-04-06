package main

import (
	"net"
)

func main() {
	conn, _ := net.Dial("udp", "127.0.0.1:8080")
	defer conn.Close()

	conn.Write([]byte("TEST"))
	conn.Write([]byte("Hello"))
	conn.Write([]byte("TEST"))
	conn.Write([]byte("TEST"))
	conn.Write([]byte("TEST"))
	conn.Write([]byte("TEST"))
	conn.Write([]byte("TEST"))

}
