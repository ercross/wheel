package main

import (
	"fmt"
	"net"
	"os"
)

const (
	serverAddr string = "127.0.0.1:15001"
)

func main() {
	resolvedServerAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Error resolving server address: ", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, resolvedServerAddr)
	if err != nil {
		fmt.Println("Error connecting to server: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Hello from the other side"))
	if err != nil {
		fmt.Println("Error writing to server: ", err)
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error reading from server: ", err)
		os.Exit(1)
	}

	fmt.Printf("Received %d bytes from server: %s\n", n, buf[:n])
	_, _ = conn.Write([]byte("exit")) // shutdown server command
}
