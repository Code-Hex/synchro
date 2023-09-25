package constraints

// Bytes is an interface that represents a byte slice.
type Bytes interface {
	[]byte | ~string
}
