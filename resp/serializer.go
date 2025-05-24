package resp

// Serialize converts a Value to its RESP byte representation.
func Serialize(v Value) []byte {
	return v.Serialize()
}
