package backendApi

import (
	DB "Renew/bin/main/aws/database"
	"Renew/bin/main/misc"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

var (
	SmsAcc string
	SmsPwd string

	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	verifyM = make(map[string]VerifyInfo)
)

func init() {
	SmsAcc = misc.Config.SmsAcc
	SmsPwd = misc.Config.SmsPwd
}

func GetVerifyCode(_len int) string {

	b := make([]byte, _len)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	fmt.Println(string(b))
	return string(b)
}

func SendVerifyCode(phone, code string) {
	cc := `歡迎註冊itLive服務，您的驗證碼為：` + code

	sms := "http://api.twsms.com/json/sms_send.php" +
		"?username=" + SmsAcc +
		"&password=" + SmsPwd +
		"&mobile=" + phone +
		"&message=" + url.QueryEscape(cc)
	res, err := http.Get(sms)
	if err != nil {
		log.Fatal(err)
	}

	//http.PostForm("", url.Values{})
	//defer res.Body.Close()

	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%s", sitemap)
	s := string(sitemap)

	fmt.Println(s)
}

func PostPhoneV2(c *gin.Context, claims *Claims) interface{} {
	type pr struct {
		Phone string `json:"phone"`
	}

	var ps pr
	err := c.ShouldBindJSON(&ps)
	if err != nil {
		return nil
	}

	vc := GetVerifyCode(6)

	verifyM[claims.UserID] = VerifyInfo{ps.Phone, vc, time.Now().Unix()}

	SendVerifyCode(ps.Phone, vc)

	//return gin.H{"state": "ok", "exp": verifyM[claims.UserID].ExpireTime + 300}

	return gin.H{"state": "ok", "exp": verifyM[claims.UserID].ExpireTime + 300, "code": vc}

}
func PostCodeV2(c *gin.Context, claims *Claims) interface{} {
	type pr struct {
		Vcode string `json:"vcode"`
	}
	var ps pr
	err := c.ShouldBindJSON(&ps)
	if err != nil {
		return nil
	}

	fmt.Println(ps)
	fmt.Println(verifyM[claims.UserID])

	var t = time.Now().Unix()

	if (t - verifyM[claims.UserID].ExpireTime) >= 300 {
		return "timeout"
	}

	if verifyM[claims.UserID].Code == ps.Vcode {
		defer delete(verifyM, claims.UserID)

		{
			km := map[string]*string{
				"#phone":  aws.String("phone"),
				"#verify": aws.String("verify"),
			}

			type updv struct {
				Phone  string `json:":phone"`
				Verify bool   `json:":verify"`
			}

			updcmd := "set #phone = :phone , #verify = :verify"

			DB.DbUpdate(DB.UserTable, "userID", claims.UserID, updcmd, km, updv{verifyM[claims.UserID].Phone, true})

		}

		return "ok"

	} else {
		return "retry"
	}
	return ""

}

func PostPhone(c *gin.Context) {

	buf := make([]byte, 128)
	n, _ := c.Request.Body.Read(buf)
	defer c.Request.Body.Close()
	s := string(buf[0:n])
	type pr struct {
		Phone string `json:"phone"`
	}

	var r pr
	err := json.Unmarshal([]byte(s), &r)
	if err != nil {
		c.String(http.StatusOK, "json error")
		fmt.Println(err)
	}

	vc := GetVerifyCode(6)

	gid, err := c.Cookie("gid")
	if err != nil {
		c.String(http.StatusOK, "cookie error")
		fmt.Println(err)
		return
	}

	//發送簡訊
	//sendVerifyCode(r.Phone, vc)
	verifyM[gid] = VerifyInfo{r.Phone, vc, time.Now().Unix()}

	//c.String(http.StatusOK, "ok")

}

func PostCode(c *gin.Context) {

	buf := make([]byte, 128)
	n, _ := c.Request.Body.Read(buf)
	defer c.Request.Body.Close()
	s := string(buf[0:n])
	fmt.Println(s)

	type verifyR struct {
		Vcode string `json:"vcode"`
	}
	var t = time.Now().Unix()
	var v verifyR

	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		c.String(http.StatusOK, "json error")
		fmt.Println(err)
		return
	}

	gid, err := c.Cookie("gid")
	if err != nil {
		c.String(http.StatusOK, "cookie error")
		fmt.Println(err)
		return
	}

	if (t - verifyM[gid].ExpireTime) >= 300 {
		c.String(http.StatusOK, "timeout")
		return
	}

	if verifyM[gid].Code == v.Vcode {
		defer delete(verifyM, gid)
		c.String(http.StatusOK, "ok")

	} else {
		c.String(http.StatusOK, "retry")
	}
}
