package handlers

import (
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

// Load the JWT key from environment
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr, err := c.Cookie("token")
        if err != nil {
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        claims := &Claims{}
        tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !tkn.Valid {
            c.Redirect(http.StatusFound, "/login")
            c.Abort()
            return
        }

        c.Set("userID", claims.UserID)
        c.Next()
    }
}

func GenerateJWT(userID int) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    return tokenString, err
}
