package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const tcpAddress = "127.0.0.1:8080"

func main() {
	listner, err := net.Listen("tcp", tcpAddress)
	if err != nil {
		fmt.Println("Error Listening", err)
		os.Exit(1)
	}
	defer listner.Close()

	fmt.Print("Listening on", tcpAddress)

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			continue
		}

		go handleConnections(conn)
	}
}

func handleConnections(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading", err)
			return
		}
		fmt.Print(message)
	}
}
