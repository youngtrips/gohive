package sms

import (
	"testing"
)

func TestSendCode(t *testing.T) {
	err := SendCode("15928010910", "123987")
	t.Logf("error: %s", err)
}
