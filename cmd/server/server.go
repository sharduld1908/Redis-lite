package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"redis-lite/kvstore"
	"redis-lite/resp"
	"strings"
	"sync"
)

type RedisServer struct {
	data  *kvstore.HashTable
	mutex sync.RWMutex
}

func NewRedisServer() *RedisServer {
	return &RedisServer{
		data:  kvstore.NewHashTable(),
		mutex: sync.RWMutex{},
	}
}

func (rs *RedisServer) sendError(writer *bufio.Writer, errorStr string) {
	errResp := resp.Error{Value: errorStr}
	_, err := writer.Write(resp.Serialize(errResp))
	if err != nil {
		log.Printf("Error sending error: %v", err)
		return
	}
	writer.Flush()
}

func (rs *RedisServer) handlePing(writer *bufio.Writer) {
	pong := resp.SimpleString{Value: "PONG"}
	_, err := writer.Write(resp.Serialize(pong))
	if err != nil {
		log.Printf("Error sending pong: %v", err)
		return
	}
}

func (rs *RedisServer) handleEcho(writer *bufio.Writer, parts []string) {
	var echo resp.BulkString
	if len(parts) > 1 {
		echo = resp.BulkString{Value: strings.Join(parts[1:], " ")}
	} else {
		echo = resp.BulkString{Value: ""}
	}
	_, err := writer.Write(resp.Serialize(echo))
	if err != nil {
		log.Printf("Error sending pong: %v", err)
		return
	}
}

func (rs *RedisServer) handleGetCommand(writer *bufio.Writer, parts []string) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	if len(parts) != 2 {
		rs.sendError(writer, "ERR wrong number of arguments for 'get' command")
		return
	}

	value, ok := rs.data.Get(parts[1])
	if !ok {
		serverResp := resp.BulkString{IsNull: true}
		_, err := writer.Write(resp.Serialize(serverResp))
		if err != nil {
			log.Printf("Error sending OK response")
			return
		}
	}

	serverResp := resp.BulkString{Value: value.(string), IsNull: false}
	_, err := writer.Write(resp.Serialize(serverResp))
	if err != nil {
		log.Printf("Error sending OK response")
		return
	}
}

func (rs *RedisServer) handleSetCommand(writer *bufio.Writer, parts []string) {
	

	if len(parts) != 3 {
		rs.sendError(writer, "ERR syntax error")
		return
	}

	rs.mutex.Lock()	
	rs.data.Insert(parts[1], parts[2])
	rs.mutex.Unlock()

	
	serverResp := resp.SimpleString{Value: "OK"}
	_, err := writer.Write(resp.Serialize(serverResp))
	if err != nil {
		log.Printf("Error sending OK response")
		return
	}
}

func (rs *RedisServer) handleHelp(writer *bufio.Writer) {
	help := resp.BulkString{Value: "PING: Returns PONG\nECHO <message>: Returns the provided message\nGET <key>: Returns the value associated with the key\nSET <key> <value>: Sets the value for the given key"}
	_, err := writer.Write(resp.Serialize(help))
	if err != nil {
		log.Printf("Error sending help: %v", err)
		return
	}
}

func (rs *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Accepted connection from %s", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		input, err := resp.Deserialize(reader)
		if err != nil {
			log.Printf("Connection closed: %v", err)
			return
		}

		log.Printf("Received from client %s", input)

		// array length cannot be zero. There has to be something in there right?
		clientArray, ok := input.(resp.Array)
		if !ok || len(clientArray.Values) == 0 {
			errResp := resp.Error{Value: "ERR invalid command format"}
			_, err = writer.Write(resp.Serialize(errResp))
			if err != nil {
				log.Printf("Error sending error: %v", err)
				return
			}
			writer.Flush()
			continue
		}

		// Get the command
		command, ok := clientArray.Values[0].(resp.BulkString)
		if !ok {
			errResp := resp.Error{Value: "ERR invalid command format"}
			_, err = writer.Write(resp.Serialize(errResp))
			if err != nil {
				log.Printf("Error sending error: %v", err)
				return
			}
			writer.Flush()
			continue
		}

		commandStr := strings.ToUpper(command.Value)
		parts := make([]string, len(clientArray.Values))
		for i, val := range clientArray.Values {
			if bulk, ok := val.(resp.BulkString); ok {
				parts[i] = bulk.Value
			} else {
				parts[i] = ""
			}
		}

		switch commandStr {
		case "PING":
			rs.handlePing(writer)
		case "ECHO":
			rs.handleEcho(writer, parts)
		case "GET":
			rs.handleGetCommand(writer, parts)
		case "SET":
			rs.handleSetCommand(writer, parts)
		case "HELP":
			rs.handleHelp(writer)
		default:
			errResp := resp.Error{Value: fmt.Sprintf("ERR unknown command '%s'", commandStr)}
			_, err = writer.Write(resp.Serialize(errResp))
			if err != nil {
				log.Printf("Error sending error: %v", err)
				return
			}
		}

		writer.Flush()
	}
}
