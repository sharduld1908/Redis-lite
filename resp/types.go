package resp

import (
	"strconv"
	"strings"
)

type Value interface {
	Serialize() []byte
	String() string
}

/*Simple String Type*/
type SimpleString struct {
	Value string
}

func (s SimpleString) Serialize() []byte {
	return []byte("+" + s.Value + "\r\n")
}

func (s SimpleString) String() string {
	return s.Value
}

/*Error Type*/
type Error struct {
	Value string
}

func (e Error) Serialize() []byte {
	return []byte("-" + e.Value + "\r\n")
}

func (e Error) String() string {
	return e.Value
}

/*Integer Type*/
type Integer struct {
	Value int64
}

func (i Integer) Serialize() []byte {
	return []byte(":" + strconv.FormatInt(i.Value, 10) + "\r\n")
}

func (i Integer) String() string {
	return strconv.FormatInt(i.Value, 10)
}

type BulkString struct {
	Value  string
	IsNull bool
}

func (b BulkString) Serialize() []byte {
	if b.IsNull {
		return []byte("$-1\r\n")
	}

	// For non-null bulk strings, format is: $<length>\r\n<string>\r\n
	length := strconv.Itoa(len(b.Value))
	return []byte("$" + length + "\r\n" + b.Value + "\r\n")
}

// String returns a human-readable representation of BulkString
func (b BulkString) String() string {
	if b.IsNull {
		return "(nil)"
	}
	return b.Value
}

type Array struct {
	Values []Value
	IsNull bool
}

func (a Array) Serialize() []byte {
	if a.IsNull {
		return []byte("*-1\r\n")
	}

	var builder strings.Builder
	builder.WriteString("*" + strconv.Itoa(len(a.Values)) + "\r\n")
	for _, v := range a.Values {
		builder.Write(v.Serialize())
	}

	return []byte(builder.String())
}

func (a Array) String() string {
	if a.IsNull {
		return "(nil)"
	}

	var elements []string
	for _, v := range a.Values {
		elements = append(elements, v.String())
	}

	return "[" + strings.Join(elements, ",") + "]"
}
