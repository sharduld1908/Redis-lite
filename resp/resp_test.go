package resp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSimpleString_Serialize(t *testing.T) {
	s := SimpleString{Value: "OK"}
	expected := []byte("+OK\r\n")
	result := s.Serialize()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Serialize() = %v, want %v", result, expected)
	}
}

func TestSimpleString_String(t *testing.T) {
	s := SimpleString{Value: "OK"}
	expected := "OK"
	result := s.String()

	if result != expected {
		t.Errorf("String() = %v, want %v", result, expected)
	}
}

func TestError_Serialize(t *testing.T) {
	e := Error{Value: "ERR something went wrong"}
	expected := []byte("-ERR something went wrong\r\n")
	result := e.Serialize()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Serialize() = %v, want %v", result, expected)
	}
}

func TestError_String(t *testing.T) {
	e := Error{Value: "ERR something went wrong"}
	expected := "ERR something went wrong"
	result := e.String()

	if result != expected {
		t.Errorf("String() = %v, want %v", result, expected)
	}
}

func TestInteger_Serialize(t *testing.T) {
	i := Integer{Value: 12345}
	expected := []byte(":12345\r\n")
	result := i.Serialize()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Serialize() = %v, want %v", result, expected)
	}
}

func TestInteger_String(t *testing.T) {
	i := Integer{Value: 12345}
	expected := "12345"
	result := i.String()

	if result != expected {
		t.Errorf("String() = %v, want %v", result, expected)
	}
}

func TestBulkString_Serialize(t *testing.T) {
	b := BulkString{Value: "hello", IsNull: false}
	expected := []byte("$5\r\nhello\r\n")
	result := b.Serialize()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Serialize() = %v, want %v", result, expected)
	}

	b = BulkString{IsNull: true}
	expected = []byte("$-1\r\n")
	result = b.Serialize()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Serialize() (null) = %v, want %v", result, expected)
	}
}

func TestBulkString_String(t *testing.T) {
	b := BulkString{Value: "hello", IsNull: false}
	expected := "hello"
	result := b.String()

	if result != expected {
		t.Errorf("String() = %v, want %v", result, expected)
	}

	b = BulkString{IsNull: true}
	expected = "(nil)"
	result = b.String()

	if result != expected {
		t.Errorf("String() (null) = %v, want %v", result, expected)
	}
}

func TestArray_Serialize(t *testing.T) {
	a := Array{
		Values: []Value{
			Integer{Value: 1},
			BulkString{Value: "two", IsNull: false},
			SimpleString{Value: "three"},
		},
		IsNull: false,
	}
	expected := []byte("*3\r\n:1\r\n$3\r\ntwo\r\n+three\r\n")
	result := a.Serialize()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Serialize() = %v, want %v", result, expected)
	}

	a = Array{IsNull: true}
	expected = []byte("*-1\r\n")
	result = a.Serialize()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Serialize() (null) = %v, want %v", result, expected)
	}
}

func TestArray_String(t *testing.T) {
	a := Array{
		Values: []Value{
			Integer{Value: 1},
			BulkString{Value: "two", IsNull: false},
			SimpleString{Value: "three"},
		},
		IsNull: false,
	}
	expected := "[1,two,three]"
	result := a.String()

	if result != expected {
		t.Errorf("String() = %v, want %v", result, expected)
	}

	a = Array{IsNull: true}
	expected = "(nil)"
	result = a.String()

	if result != expected {
		t.Errorf("String() (null) = %v, want %v", result, expected)
	}
}

func TestSerialize(t *testing.T) {
	i := Integer{Value: 123}
	expected := []byte(":123\r\n")
	result := Serialize(i)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Serialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_SimpleString(t *testing.T) {
	input := []byte("+OK\r\n")
	reader := bytes.NewReader(input)
	expected := SimpleString{Value: "OK"}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_Error(t *testing.T) {
	input := []byte("-ERR something\r\n")
	reader := bytes.NewReader(input)
	expected := Error{Value: "ERR something"}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_Integer(t *testing.T) {
	input := []byte(":123\r\n")
	reader := bytes.NewReader(input)
	expected := Integer{Value: 123}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_BulkString(t *testing.T) {
	input := []byte("$5\r\nhello\r\n")
	reader := bytes.NewReader(input)
	expected := BulkString{Value: "hello"}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}

	input = []byte("$-1\r\n")
	reader = bytes.NewReader(input)
	expected = BulkString{IsNull: true}

	result, err = Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() (null) = %v, want %v", result, expected)
	}
}

func TestDeserialize_Array(t *testing.T) {
	input := []byte("*3\r\n:1\r\n$3\r\ntwo\r\n+three\r\n")
	reader := bytes.NewReader(input)
	expected := Array{
		Values: []Value{
			Integer{Value: 1},
			BulkString{Value: "two"},
			SimpleString{Value: "three"},
		},
	}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}

	input = []byte("*-1\r\n")
	reader = bytes.NewReader(input)
	expected = Array{IsNull: true}
	result, err = Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() (null) = %v, want %v", result, expected)
	}
}

func TestDeserialize_InvalidType(t *testing.T) {
	input := []byte("%invalid\r\n")
	reader := bytes.NewReader(input)

	_, err := Deserialize(reader)
	if err == nil {
		t.Errorf("Deserialize() should return an error for invalid type")
	}
}

func TestDeserialize_BulkString_Empty(t *testing.T) {
	input := []byte("$0\r\n\r\n")
	reader := bytes.NewReader(input)
	expected := BulkString{Value: ""}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_Array_SingleBulkString(t *testing.T) {
	input := []byte("*1\r\n$4\r\nping\r\n")
	reader := bytes.NewReader(input)
	expected := Array{
		Values: []Value{
			BulkString{Value: "ping"},
		},
	}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_Array_TwoBulkStrings(t *testing.T) {
	input := []byte("*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n")
	reader := bytes.NewReader(input)
	expected := Array{
		Values: []Value{
			BulkString{Value: "echo"},
			BulkString{Value: "hello world"},
		},
	}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_Array_GetKeyValue(t *testing.T) {
	input := []byte("*2\r\n$3\r\nget\r\n$3\r\nkey\r\n")
	reader := bytes.NewReader(input)
	expected := Array{
		Values: []Value{
			BulkString{Value: "get"},
			BulkString{Value: "key"},
		},
	}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_SimpleString_HelloWorld(t *testing.T) {
	input := []byte("+hello world\r\n")
	reader := bytes.NewReader(input)
	expected := SimpleString{Value: "hello world"}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_Error_ErrorMessage(t *testing.T) {
	input := []byte("-Error message\r\n")
	reader := bytes.NewReader(input)
	expected := Error{Value: "Error message"}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}

func TestDeserialize_BulkString_Null(t *testing.T) {
	input := []byte("$-1\r\n")
	reader := bytes.NewReader(input)
	expected := BulkString{IsNull: true}

	result, err := Deserialize(reader)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Deserialize() = %v, want %v", result, expected)
	}
}
