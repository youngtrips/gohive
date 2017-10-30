package registry

import (
	"testing"
)

func TestRegistry(t *testing.T) {
	c := One("auth")
	if c == nil {
		t.Logf("no found")
	} else {
		t.Logf(c.Addr)
	}
}
