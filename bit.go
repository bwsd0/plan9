package plan9

/*
The gbit and pbit helper functions that translate and pack numbers and byte
sequences in little-endian byte order.
*/

// FIXME(bwsd): bounds checking elimination is not performed.
// See: golang.org/issue/14808

// gbit8 decodes a uint8 from b and returns that value and the remaining slice
// of b.
func gbit8(b []byte) (uint8, []byte) {
	return uint8(b[0]), b[1:]
}

// gbit16 decodes a uint16 from b and returns that value and the remaining slice
// of b.
func gbit16(b []byte) (uint16, []byte) {
	// _ = b[1] compiler BCE hint; see: golang.org/issue/14808
	return uint16(b[0]) | uint16(b[1])<<8, b[2:]
}

// gbit32 decodes a uint32 from from b and returns that value and the remaining
// slice of b.
func gbit32(b []byte) (uint32, []byte) {
	// _ = b[3] compiler BCE hint; see: golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, b[4:]
}

// gbit64 decodes a uint64 from b and returns that value and and the remaining
// slice of b.
func gbit64(b []byte) (uint64, []byte) {
	lo, b := gbit32(b)
	hi, b := gbit32(b)
	return uint64(hi)<<32 | uint64(lo), b
}

// gstring reads a string from b, prefixed with a 16-bit length. It returns the
// string and the remaining slice of b.
func gstring(b []byte) (string, []byte) {
	n, b := gbit16(b)
	return string(b[0:n]), b[n:]
}

// pbit8 encodes a uint8 x into b and returns the remaining slice of b.
func pbit8(b []byte, x uint8) []byte {
	n := len(b)
	if n+1 > cap(b) {
		nb := make([]byte, n, 100+2*cap(b))
		copy(nb, b)
		b = nb
	}
	b = b[0 : n+1]
	b[n] = x
	return b
}

// pbit16 encodes a uint16 x into b and returns the remaining slice of b.
func pbit16(b []byte, x uint16) []byte {
	n := len(b)
	if n+2 > cap(b) {
		nb := make([]byte, n, 100+2*cap(b))
		copy(nb, b)
		b = nb
	}
	b = b[0 : n+2]
	b[n] = byte(x)
	b[n+1] = byte(x >> 8)
	return b
}

// pbit32 encodes the uint32 x into b and returns the remaining slice of b.
func pbit32(b []byte, x uint32) []byte {
	n := len(b)
	if n+4 > cap(b) {
		nb := make([]byte, n, 100+2*cap(b))
		copy(nb, b)
		b = nb
	}
	b = b[0 : n+4]
	b[n] = byte(x)
	b[n+1] = byte(x >> 8)
	b[n+2] = byte(x >> 16)
	b[n+3] = byte(x >> 24)
	return b
}

// pbit64 encodes the uint64 x into b and returns the remaining slice of b.
func pbit64(b []byte, x uint64) []byte {
	b = pbit32(b, uint32(x))
	b = pbit32(b, uint32(x>>32))
	return b
}

// pstring copies the string s to b, prepending it with a 16-bit length and
// returning the remaining slice of b.
//
// If the buffer is too small, pstring will panic.
func pstring(b []byte, s string) []byte {
	if len(s) >= 1<<16 {
		panic(ProtocolError("string too long"))
	}
	b = pbit16(b, uint16(len(s)))
	b = append(b, []byte(s)...)
	return b
}
