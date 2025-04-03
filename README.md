# Redis Lite

A lightweight Redis server implementation in Go that supports basic Redis commands and the RESP (Redis Serialization Protocol) for communication.

## Overview

Redis Lite is a simplified Redis server that implements the core functionality of Redis including:

- RESP (Redis Serialization Protocol) serialization and deserialization
- Basic key-value operations
- Simple string, error, integer, bulk string, and array data types
- TCP server implementation with concurrent client handling

## Getting Started

### Prerequisites

- Go (1.16 or later recommended)
- Make

### Building and Running

```bash
# Build the project
make

# Start the Redis Lite server
make run-server

# Connect a client to the server
make run-client
```

### Using with Redis CLI

The server is compatible with standard Redis clients. You can connect using:

```bash
# Connect with redis-cli
redis-cli -h localhost -p 5000

# Run benchmark tests
redis-benchmark -h localhost -p 5000
```

## Code Structure

- **resp package**: Implements Redis Serialization Protocol
  - `types.go`: Defines the RESP data types (SimpleString, Error, Integer, BulkString, Array)
  - `serializer.go`: Converts RESP values to byte representation
  - `deserializer.go`: Parses RESP protocol data from byte streams

- **main package**: Implements the server
  - `main.go`: Entry point that starts TCP server on port 5000
  - `server.go`: Handles client connections and implements Redis commands

## Supported Commands

The server currently supports the following commands:

- `PING`: Returns PONG
- `ECHO <message>`: Returns the provided message
- `GET <key>`: Retrieves the value associated with the specified key
- `SET <key> <value>`: Stores a value with the specified key
- `HELP`: Shows available commands and their usage

## Technical Implementation

### RESP Protocol

The implementation includes a full RESP protocol parser and serializer that handles:
- Simple strings (Status replies): `+<string>\r\n`
- Errors: `-<error message>\r\n`
- Integers: `:<number>\r\n`
- Bulk strings: `$<length>\r\n<data>\r\n`
- Arrays: `*<count>\r\n<elements...>`

### Concurrency

The server uses Go's goroutines to handle multiple client connections concurrently and a mutex to ensure thread-safe access to the shared data store.

### Data Storage

A simple in-memory map is used to store key-value pairs, guarded by a mutex for thread safety when handling concurrent client requests.

## Future Improvements

Potential enhancements:
- Add more Redis commands (DEL, EXISTS, etc.)
- Implement data persistence
- Add TTL (Time To Live) support for keys
- Support for more complex data structures (lists, sets, sorted sets)