package main

import (
	"testing"
)

// Her fonksiyonumuza test yazmalıyız. Test ismiyle başlamalıyız.
// bu dizine geldikten sonra > go test -v
func TestRun(t *testing.T) {
	_, err := run()
	if err != nil {
		t.Error("failed run()")
	}
}
