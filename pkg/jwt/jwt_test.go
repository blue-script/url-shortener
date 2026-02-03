package jwt_test

import (
	"testing"

	"github.com/blue-script/url-shortener/pkg/jwt"
)

func TestJWTCreate(t *testing.T) {
	const email = "test@test.com"
	jwtService := jwt.NewJWT("H7TW58QykamloVtXelhA+ihasqyacBjeI4ahNFajhDI=")
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is invalid")
	}
	if data.Email != email {
		t.Fatalf("Email %s not equal %s", data.Email, email)
	}
}
