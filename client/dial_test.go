package client

import (
	"testing"
)

func TestNamespace(t *testing.T) {
	ns := `/tmp/ns/` + getuser() + `.:0`
	nsTests := []struct {
		disp string
		want string
	}{
		{":0", ns},
		{":0.1234", ns},
		{":0.0", ns},
		{"hostfoo:0", ns},
		{":0.0", ns},
	}
	for _, tt := range nsTests {
		if got := Namespace(); got != tt.want {
			t.Errorf("NameSpace(%q) == %q; want %q", tt.disp, got, tt.want)
		}
	}
}
