package backendApi

import "github.com/dgrijalva/jwt-go"

type VerifyInfo struct {
	Phone      string
	Code       string
	ExpireTime int64
}
type GoogleUserInfo struct {
	GoogleID  string `json:"id"`
	Email     string `json:"email"`
	NameFirst string `json:"given_name"`
	NameLast  string `json:"family_name"`
	Picture   string `json:"picture"`
}

type Claims struct {
	UserID     string `json:"userID"`
	UserName   string `json:"userName"`
	IsSeller   bool   `json:"isSeller"`
	ChannelID  string `json:"channelId"`
	ChannelUrl string `json:"channelUrl"`
	jwt.StandardClaims
}
