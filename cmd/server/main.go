package main

import (
	"log"
	"net"
)

const (
	listenAddr = "localhost:5000"
)

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start the Redis Lite Server at port:%s. Error: %v", listenAddr, err)
	}

	defer listener.Close()

	log.Printf("TCP Server started. Listening on %s", listenAddr)
	log.Println("Waiting for clients to connect...")

	rs := NewRedisServer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go rs.handleConnection(conn)
	}

}
