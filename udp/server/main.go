package main

import (
	"fmt"
	"net"
	"os"
)

const (
	port int    = 15001
	ip   string = "127.0.0.1"
)

func main() {
	addr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("Listening on UDP port %d\n", port)

	for {
		buf := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP: " + err.Error())
			continue
		}
		if clientAddr.Port == 0 || clientAddr.IP.Equal(net.ParseIP("0.0.0.0")) {
			continue
		}
		message := string(buf[:n])
		fmt.Printf("Received data (%s) from UDP address %s\n", message, clientAddr.String())

		if message == "exit" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}

		// echo back to client address
		_, err = conn.WriteToUDP(buf[:n], clientAddr)
		if err != nil {
			fmt.Println("Error writing to UDP: " + err.Error())
			continue
		}
	}
}
