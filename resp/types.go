package resp

import (
	"bytes"
	"strconv"
	"strings"
)

// Predefined CRLF constant
var crlf = []byte("\r\n")

type Value interface {
	Serialize() []byte
	String() string
}

/*Simple String Type*/
type SimpleString struct {
	Value string
}

func (s SimpleString) Serialize() []byte {
	buf := make([]byte, 0, 1+len(s.Value)+2)
	buf = append(buf, '+')
	buf = append(buf, s.Value...)
	buf = append(buf, crlf...)
	return buf
}

func (s SimpleString) String() string {
	return s.Value
}

/*Error Type*/
type Error struct {
	Value string
}

func (e Error) Serialize() []byte {
	buf := make([]byte, 0, 1+len(e.Value)+2)
	buf = append(buf, '-')
	buf = append(buf, e.Value...)
	buf = append(buf, crlf...)
	return buf
}

func (e Error) String() string {
	return e.Value
}

/*Integer Type*/
type Integer struct {
	Value int64
}

func (i Integer) Serialize() []byte {
	buf := make([]byte, 0, 23)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, i.Value, 10)
	buf = append(buf, crlf...)
	return buf
}

func (i Integer) String() string {
	return strconv.FormatInt(i.Value, 10)
}

type BulkString struct {
	Value  string
	IsNull bool
}

// Predefined Null Bulk String constant
var nullBulkBytes = []byte("$-1\r\n")
var emptyBulkBytes = []byte("$0\r\n\r\n") // Optimization for empty string

// Optimized BulkString
func (b BulkString) Serialize() []byte {
	if b.IsNull {
		return nullBulkBytes
	}
	if len(b.Value) == 0 {
		return emptyBulkBytes // Specific optimization for empty string
	}

	// Format: $<length>\r\n<value>\r\n
	lenStr := strconv.Itoa(len(b.Value)) // Still need string length
	// Estimate size: 1 ('$') + len(lenStr) + 2 + len(value) + 2
	buf := make([]byte, 0, 1+len(lenStr)+2+len(b.Value)+2)

	buf = append(buf, '$')
	buf = append(buf, lenStr...)
	buf = append(buf, crlf...)
	buf = append(buf, b.Value...)
	buf = append(buf, crlf...)
	return buf
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

// Predefined Null Array constant
var nullArrayBytes = []byte("*-1\r\n")

// Optimized Array
func (a Array) Serialize() []byte {
	if a.IsNull {
		return nullArrayBytes
	}

	// Use bytes.Buffer as element serialization complexity varies
	var buf bytes.Buffer

	buf.WriteByte('*')
	buf.WriteString(strconv.Itoa(len(a.Values))) // WriteString is efficient on buffer
	buf.Write(crlf)                              // Write CRLF bytes

	for _, v := range a.Values {
		// Assume sub-serializers are also optimized
		buf.Write(v.Serialize()) // Write serialized bytes of element
	}
	return buf.Bytes() // .Bytes() avoids extra copy if possible
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
