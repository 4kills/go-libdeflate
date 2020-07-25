package native

func startMemAddress(b []byte) *byte {
	if len(b) > 0 {
		return &b[0]
	}

	b = append(b, 0)
	ptr := &b[0]
	b = b[0:0]

	return ptr
}
