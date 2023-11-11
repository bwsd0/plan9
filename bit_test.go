package plan9

import (
	"strings"
	"testing"
)

func TestPstringTooLong(t *testing.T) {
	tt := struct {
		name string
		b    []byte
		s    string
	}{
		name: "string too long",
		b:    make([]byte, STATMAX),
		s:    strings.Repeat("b", STATMAX+1),
	}
	t.Run(tt.name, func(t *testing.T) {
		defer func() { _ = recover() }()
		_ = pstring(tt.b, tt.s)
		t.Errorf("expecting panic for %s", tt.name)
	})
}

var pbuf = []byte{0, 0, 0, 0, 0, 0, 0, 0}

func BenchmarkPBit8(b *testing.B) {
	i := 0
	b.SetBytes(2)
	for ; i < b.N; i++ {
		_ = pbit8(pbuf[:2], uint8(i))
	}
}

func BenchmarkPBit16(b *testing.B) {
	i := 0
	b.SetBytes(2)
	for ; i < b.N; i++ {
		_ = pbit16(pbuf[:2], uint16(i))
	}
}

func BenchmarkPBit32(b *testing.B) {
	i := 0
	b.SetBytes(4)
	for ; i < b.N; i++ {
		_ = pbit32(pbuf[:4], uint32(i))
	}
}

func BenchmarkPBit64(b *testing.B) {
	i := 0
	b.SetBytes(8)
	for ; i < b.N; i++ {
		_ = pbit64(pbuf[:8], uint64(i))
	}
}
