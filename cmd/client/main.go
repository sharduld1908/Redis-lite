package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"redis-lite/resp"
	"strings"
)

const (
	serverAddr = "localhost:5000"
)

func main() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect to server at %s: %v", serverAddr, err)
	}
	// Ensure the connection is closed when main exits
	defer conn.Close()

	log.Printf("Connected to server: %s", serverAddr)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	consoleReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := consoleReader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.Split(input, " ")

		var respArray resp.Array
		respArray.Values = make([]resp.Value, len(parts))

		for i, part := range parts {
			respArray.Values[i] = resp.BulkString{Value: part}
		}

		_, err = writer.Write(resp.Serialize(respArray))
		if err != nil {
			log.Printf("Error sending command: %v", err)
			return
		}
		writer.Flush()

		serverResp, err := resp.Deserialize(reader)
		if err != nil {
			log.Printf("Error reading server response: %v", err)
			return
		}

		fmt.Println(serverResp.String())
	}
}
