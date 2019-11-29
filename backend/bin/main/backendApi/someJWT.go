package backendApi

import (
	"Renew/bin/main/misc"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	TokenAccecpt = iota
	TokenTimeout
	TokenValidationFailed
)

var jwtSecret = []byte(misc.Config.JwtSecret)

func GenerateToken(userID, UserName, ChannelID, ChannelUrl string, isSeller bool) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		userID,
		UserName,
		isSeller,
		ChannelID,
		ChannelUrl,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "renewBB",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	if debugMode {
		fmt.Println(claims)
		fmt.Println(token)
	}

	return token, err
}

func ParseToken(token string) (int, *Claims, error) {

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			if time.Now().Unix() > claims.ExpiresAt {
				return TokenTimeout, nil, nil
			}
			return TokenAccecpt, claims, nil
		}
	}
	return TokenValidationFailed, nil, err
}

// Jwt example
func VerifyJWT(c *gin.Context) {

	var state string
	type resp struct {
		FullToken        string      `json:"full_token"`
		DecryptedPayload interface{} `json:"decrypted_payload"`
	}
	var data resp

	token, err := c.Cookie("jwtAccess")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(token)

	if token == "" {
		state = "notLogin"
		c.JSONP(http.StatusOK, gin.H{
			"state": state,
		})
		return
	} else {
		tokenState, claims, err := ParseToken(token)
		if err != nil {
			fmt.Println(err)
		}
		switch tokenState {
		case TokenAccecpt:
			state = "hadLogin"
			data.FullToken = token
			data.DecryptedPayload = claims
		case TokenTimeout:
			state = "timeout"
		case TokenValidationFailed:
			state = "tokenFail"
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"state": state,
		"data":  data,
	})

	if debugMode {
		fmt.Println(state)
		fmt.Println(data)
	}

}
