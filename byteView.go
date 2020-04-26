package birdcache

// ByteView holds an imutable view of bytes
type ByteView struct {
	b []byte
}

// Len returns the view's size
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice return the copy of the view as a byte slice
// because it is read only
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String return the data to human language
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
