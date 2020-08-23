package gocache

// ByteView holds immutable view of bytes
type ByteView struct {
	b []byte
}

// Len returns the length of the byteview
func (v ByteView) Len() int {
	if v.b == nil {
		return 0
	}
	return len(v.b)
}

// ByteSlice returns a copy of the data as a slice
func (v ByteView) ByteSlice() []byte {
	if v.b == nil {
		return nil
	}
	return cloneBytes(v.b)
}

// String stringifies the bytes of the byteview
func (v ByteView) String() string {
	if v.b == nil {
		return ""
	}
	return string(v.b)
}

// cloneBytes returns a copy of the src byte slice
func cloneBytes(bs []byte) []byte {
	r := make([]byte, len(bs))
	copy(r, bs)
	return r
}
