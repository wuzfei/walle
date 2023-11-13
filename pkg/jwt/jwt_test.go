package jwt

import (
	"fmt"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	j := NewJwt(&Config{
		TokenExpire:        time.Second + 600,
		RefreshTokenExpire: time.Second + 6000,
		Key:                "afweaf",
	})
	s, ts, err := j.CreateToken(TokenPayload{
		UserId:   23,
		Email:    "fawefaq@qq.com",
		Username: "fwafwef",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s, ts, err)

	tp, err := j.ValidateToken(s)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tp.Email, err)

}
