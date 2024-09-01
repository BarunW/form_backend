package data

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaim struct{
    Email string 
    jwt.RegisteredClaims
}

type tokens struct{
    AccessToken string
    RefreshToken string 
}

func newAccessToken(email string) (string, error) {
    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaim{
        Email: email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
            Issuer: "form_poc",
            IssuedAt: jwt.NewNumericDate(time.Now()),
        },
    }).SignedString([]byte("accesstokensecret")) 

    return token, err
}

func(s *PostgresStore) setUserRefreshToken(user_id int, email string) (string,error){ 
    var(
        err error
        exist bool
        refreshToken string
    )

    refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaim{
        Email: email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
            Issuer: "form_poc",
            IssuedAt: jwt.NewNumericDate(time.Now()),
        },
    }).SignedString([]byte("refresh_token_secret")) 

    if err != nil{
        slog.Error("Unable to generate newRefresh token", "details", err.Error())
        return "", err
    }
    
    exist, err = s.isRefreshTokenExist(user_id) 
    if err != nil{
        return "", err 
    }
    if exist == false{
        fmt.Println("inserting refresh token")
        return s.insertRefreshToken(user_id , refreshToken) 
    }
    
    fmt.Println("Updating refresh token")
    return s.updateRefreshToken(user_id, refreshToken)

}

func(s *PostgresStore) updateRefreshToken(uid int, refreshToken string) (string, error) {
    var err error
    updateQuery := fmt.Sprintf("UPDATE jwt_tokens_ SET token=array_append(token, '%s') WHERE user_id=%d;",refreshToken, uid)
    fmt.Println(updateQuery)
    _, err = s.db.Exec(updateQuery)
    if err != nil{
        slog.Error("Failed to set user refresh token", "details", err.Error())
        return "", err
    }
    
    return refreshToken, nil 
}

func(s *PostgresStore) insertRefreshToken(uid int, refreshToken string) (string, error) {
    var err error

    insertQuery := fmt.Sprintf("INSERT INTO jwt_tokens_(user_id, token) VALUES(%d,'{%s}')",uid, refreshToken)
    _, err = s.db.Exec(insertQuery)
    if err != nil{
        slog.Error("Failed to set user refresh token", "details", err.Error())
        return "", err
    }
    
    return refreshToken, nil 
}

func(s *PostgresStore) isRefreshTokenExist(user_id int) (bool, error){ 
    var exist bool = true 
    err := s.db.QueryRow("SELECT EXISTS ( SELECT 1 FROM jwt_tokens_ WHERE user_id=$1);", user_id).Scan(&exist)
    if err != nil {
        slog.Error("Failed to check whether token exist or not", "details", err.Error())
        return exist, err
    }
    fmt.Println(exist)
    if exist == false{
        return false, nil 
    }

    return true, nil 
}



