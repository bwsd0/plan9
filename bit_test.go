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
		b:    make([]byte, 1<<16),
		s:    strings.Repeat("b", (1<<16)+1),
	}
	t.Run(tt.name, func(t *testing.T) {
		defer func() { _ = recover() }()
		_ = pstring(tt.b, tt.s)
		t.Errorf("expecting panic for %s", tt.name)
	})
}
