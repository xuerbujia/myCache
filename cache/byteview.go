package cache

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}
func cloneBytes(b []byte) []byte {
	res := make([]byte, len(b))
	copy(res, b)
	return res
}
func (v ByteView) String() string {
	return string(v.b)
}
