package utils

import (
	"testing"
)

func TestCreateJwtToken(t *testing.T) {
	token, err := CreateJwtToken("codelee1", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
	jwtInfo, err := ParseToken(token)
	if err != nil {
		t.Error(err)
	}
	t.Log(jwtInfo)
}
