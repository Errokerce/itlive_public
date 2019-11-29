package abandon

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	DB "Renew/bin/main/aws/database"
	api "Renew/bin/main/backendApi"
)

const (
	TokenAccecpt = iota
	TokenTimeout
	TokenValidationFailed
	debugMode = false
)

func SetChannelID(c *gin.Context, claims *api.Claims) interface{} {

	type jsonChannel struct {
		ChannelID string `json:"channel_id"`
		CustomUrl string `json:"custom_url"`
	}
	var jscid jsonChannel

	err := c.ShouldBindJSON(&jscid)
	if err != nil {
		fmt.Println(err)
	}

	jscid.ChannelID = strings.TrimPrefix(jscid.ChannelID, "https://www.youtube.com/channel/")

	filter := expression.Or(expression.Name("channelID").Equal(expression.Value(aws.String(jscid.ChannelID))),
		expression.Name("customUrl").Equal(expression.Value(aws.String(jscid.CustomUrl))))
	project := expression.NamesList(expression.Name("channelID"), expression.Name("customUrl"))

	type chk struct {
		ChannelID string `json:"channelID"`
		CustomUrl string `json:"customUrl"`
	}
	cl := DB.DbQueryMany(DB.UserTable, filter, project, chk{})

	if len(cl) > 0 {
		return gin.H{
			"msg":  "repeat",
			"data": cl,
		}
	}

	km := map[string]*string{
		"#cid": aws.String("channelID"),
		"#url": aws.String("customUrl"),
	}
	type updv struct {
		ChannelID string `json:":cid"`
		CustomUrl string `json:":url"`
	}

	cmd := "set #cid = :cid , #url = :url"
	DB.DbUpdate(DB.UserTable, "userID", claims.UserID, cmd, km, updv{jscid.ChannelID, jscid.CustomUrl})
	return "ok"
}

func PostIsSeller(c *gin.Context) {

	var state string
	var u DB.User

	token, err := c.Cookie("jwtAccess")
	if err != nil {
		fmt.Println(err)
	}

	if token == "" {
		state = "notLogin"
		c.JSONP(http.StatusOK, gin.H{
			"state": state,
		})
		return
	} else {
		tokenState, claims, err := api.ParseToken(token)
		if err != nil {
			fmt.Println(err)
		}
		switch tokenState {
		case TokenAccecpt:

			type CC struct {
				IsSeller bool `json:"is_seller"`
			}
			var s CC
			err = c.ShouldBindJSON(&s)
			fmt.Println(s)

			kn := map[string]*string{"#o": aws.String("isSeller")}

			type updv struct {
				IsSeller bool `json:":o"`
			}

			u := updv{s.IsSeller}
			DB.DbUpdate(DB.UserTable, "userID", claims.UserID,
				"set #o = :o", kn, u)

			state = "ok"
		case TokenTimeout:
			state = "timeout"
		case TokenValidationFailed:
			state = "tokenFail"
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"state": state,
		"data":  u,
	})

	if debugMode {
		fmt.Println(state)
		fmt.Println(u)
	}

}

func ViewUserdata(c *gin.Context) {

	var state string
	var u DB.User

	token, err := c.Cookie("jwtAccess")
	if err != nil {
		fmt.Println(err)
	}

	if token == "" {
		state = "notLogin"
		c.JSONP(http.StatusOK, gin.H{
			"state": state,
		})
		return
	} else {
		tokenState, claims, err := api.ParseToken(token)
		if err != nil {
			fmt.Println(err)
		}
		switch tokenState {
		case TokenAccecpt:
			DB.DbQuerybyKey(DB.UserTable, "userID", claims.UserID, &u)
			state = "ok"
		case TokenTimeout:
			state = "timeout"
		case TokenValidationFailed:
			state = "tokenFail"
		}
	}

	c.JSONP(http.StatusOK, gin.H{
		"state": state,
		"data":  u,
	})

	if debugMode {
		fmt.Println(state)
		fmt.Println(u)
	}

}
