package backendApi

import (
	"Renew/bin/main/aws/database"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

const debugMode = false
const localTest = "http://localhost:8080/login/GoogleCallback"
const onlineUrl = "https://api.itlive.nctu.me/login/GoogleCallback"
const cookieUrl = ".itlive.nctu.me"
const lvh = "localhost"
const cookieTime = 60 * 60 * 24 * 14

var (
	cookieDomain = cookieUrl
	redirectURL  = onlineUrl
)

var endpotin = oauth2.Endpoint{
	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	TokenURL: "https://accounts.google.com/o/oauth2/token",
}
var googleOauthConfig = &oauth2.Config{
	ClientID:     "",//google oauth id
	ClientSecret: "",//google oauth secret
	RedirectURL:  redirectURL,
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/youtube.readonly"},
	Endpoint: google.Endpoint,
}

const oauthStateString = "Ididit"

func HandleGoogleLogin(c *gin.Context) {
	goc := googleOauthConfig

	url := goc.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	//fmt.Println(state)

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	var r GoogleUserInfo

	if err := json.Unmarshal([]byte(contents), &r); err != nil {
		fmt.Println("ERROR:", err)
	}

	//var u database.User
	//database.DbQuerybyKey(database.UserTable, "userID", r.GoogleID, &u)

	type CheckUserExist struct {
		UserID string `json:"userID"`
		Verify bool   `json:"verify"`
	}
	filter := expression.Name("userID").Equal(expression.Value(aws.String(r.GoogleID)))
	project := expression.NamesList(expression.Name("userID"), expression.Name("verify"))

	backC := database.DbQueryMany(database.UserTable, filter, project, CheckUserExist{})
	//fmt.Printf("backC\n%v\n", backC)
	//temp := GetProfileEdit(nil, &Claims{r.GoogleID, "", false, "", "", jwt.StandardClaims{}})
	//fmt.Println(temp)
	//fmt.Println(reflect.TypeOf(temp))
	fmt.Printf("len(backC)\n%v\n", len(backC))

	if len(backC) <= 0 {
		println("newUser")
		database.DbAdd(database.UserTable, database.NewUser(r.GoogleID, r.NameLast, r.NameFirst, r.Email))
		jwtToken := FullToken(r.GoogleID)
		c.SetCookie("jwtAccess", jwtToken, cookieTime, "/", cookieDomain, false, true)

		c.Redirect(http.StatusPermanentRedirect, "https://www.itlive.nctu.me/PhoneVerify")
	} else {
		var u CheckUserExist
		PrintErr(mapstructure.Decode(backC[0], &u))
		fmt.Printf("u\n%v\n", u)
		jwtToken := FullToken(u.UserID)
		c.SetCookie("jwtAccess", jwtToken, cookieTime, "/", cookieDomain, false, true)

		if u.Verify {
			c.Redirect(http.StatusPermanentRedirect, "https://www.itlive.nctu.me")
		} else {
			c.Redirect(http.StatusPermanentRedirect, "https://www.itlive.nctu.me/PhoneVerify")
		}
	}

	//c.SetCookie("gid", string(r.GoogleID), 3000, "/", "localhost", false, false)

	//c.String(200, string(contents))

	//c.Redirect(http.StatusTemporaryRedirect, "/verify")
	if debugMode {
		//fmt.Println("Content: \n", string(contents))
		//fmt.Println("r")
		//fmt.Println(r)
		//fmt.Println("u")
		//fmt.Println(u)
	}

}
