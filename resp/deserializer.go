package resp

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

func Deserialize(reader *bufio.Reader) (Value, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch firstByte {
	case '+':
		return DeserializerSimpleString(reader)
	case '-':
		return DeserializerError(reader)
	case ':':
		return DeserializerInteger(reader)
	case '$':
		return DeserializerBulkString(reader)
	case '*':
		return DeserializerArray(reader)
	default:
		return nil, errors.New("unknown RESP data type")
	}
}

func DeserializerSimpleString(reader *bufio.Reader) (SimpleString, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return SimpleString{}, err
	}

	return SimpleString{Value: strings.TrimRight(line, "\r\n")}, nil
}

func DeserializerError(reader *bufio.Reader) (Error, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return Error{}, err
	}

	return Error{Value: strings.TrimRight(line, "\r\n")}, nil
}

func DeserializerInteger(reader *bufio.Reader) (Integer, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return Integer{}, err
	}
	i, err := strconv.ParseInt(strings.TrimRight(line, "\r\n"), 10, 64)
	if err != nil {
		return Integer{}, err
	}

	return Integer{Value: i}, nil
}

func DeserializerBulkString(reader *bufio.Reader) (BulkString, error) {
	lengthStr, err := reader.ReadString('\n')
	if err != nil {
		return BulkString{}, err
	}

	length, err := strconv.Atoi(strings.TrimRight(lengthStr, "\r\n"))
	if err != nil {
		return BulkString{}, err
	}

	if length == -1 {
		return BulkString{IsNull: true}, nil
	}

	buf := make([]byte, length+2)
	_, err = reader.Read(buf)
	if err != nil {
		return BulkString{}, err
	}
	return BulkString{Value: string(buf[:length])}, nil
}

func DeserializerArray(reader *bufio.Reader) (Array, error) {
	lengthLine, err := reader.ReadString('\n')
	if err != nil {
		return Array{}, err
	}

	length, err := strconv.Atoi(strings.TrimRight(lengthLine, "\r\n"))
	if err != nil {
		return Array{}, err
	}

	if length == -1 {
		return Array{IsNull: true}, nil
	}

	values := make([]Value, length)
	for i := range length {
		v, err := Deserialize(reader)
		if err != nil {
			return Array{}, err
		}
		values[i] = v
	}

	return Array{Values: values}, nil
}
