package main

import (
	"fmt"
	"net"
)

func main() {
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 8080})
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, addr, _ := conn.ReadFromUDP(buf)
		fmt.Printf("[%s] %x\n", addr, buf[:n])
	}
}
