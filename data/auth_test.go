package data_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaim struct{
    Email string 
    jwt.RegisteredClaims
}

func TestAccessToken(t *testing.T) {

    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaim{
        Email: "whatever@gmail.com",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
            Issuer: "form_poc",
            IssuedAt: jwt.NewNumericDate(time.Now()),
        },
    }).SignedString([]byte("accesstokensecret"))
    
    if err != nil{
        t.Fatal(err)
    }

    fmt.Print(token, err)
}
