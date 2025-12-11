package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	count := os.Getenv("COLLATZ_COUNT")
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	if count == "" || host == "" || port == "" {
		log.Fatal("COLLATZ_COUNT, SERVER_HOST, and SERVER_PORT must be set")
	}

	address := net.JoinHostPort(host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(count + "\n"))
	if err != nil {
		log.Fatalf("Failed to send data: %v", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Average steps: %s", strings.TrimSpace(response))
}
