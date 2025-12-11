package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func collatzSteps(n int) int {
	steps := 0
	for n != 1 {
		if n%2 == 0 {
			n = n / 2
		} else {
			n = 3*n + 1
		}
		steps++
	}
	return steps
}

func calculateAverage(n int) float64 {
	totalSteps := 0
	for i := 1; i <= n; i++ {
		totalSteps += collatzSteps(i)
	}
	return float64(totalSteps) / float64(n)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading: %v", err)
		return
	}

	input = strings.TrimSpace(input)
	n, err := strconv.Atoi(input)
	if err != nil || n <= 0 {
		log.Printf("Invalid number: %s", input)
		return
	}

	average := calculateAverage(n)
	response := fmt.Sprintf("%.2f\n", average)

	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error writing: %v", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Println("Server listening on port 9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Println("Client connected")
		go handleConnection(conn)
	}
}
