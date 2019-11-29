package backendApi

import (
	DB "Renew/bin/main/aws/database"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	ytKey     = ""//youtube的api token
	streamMap map[string]StreamInfo
)

type StreamInfo struct {
	Url        string
	ExpireTime time.Time
}

type Reply struct {
	Items []ytItem `json:"items,omitempty"`
}
type ytItem struct {
	Id Id `json:"id"`
	//Snippet Snippet `json:"snippet"`
}
type Id struct {
	// Kind        string  `json:"kind"`
	VideoID string `json:"videoID"`
}

//type Snippet struct {
//	// PublishedAt string  `json:"publishedAt"`
//	// ChannelId   string  `json:"channelId"`
//	Title      string     `json:"title"`
//	Thumbnails Thumbnails `json:"thumbnails"`
//	// Description string  `json:"description"`
//}
//type Thumbnail struct {
//	Url string `json:"url"`
//	//Width  string `json:"width"`
//	//Height string `json:"height"`
//}
//type Thumbnails struct {
//	// Default Thumbnail   `json:"default"`
//	High Thumbnail `json:"high"`
//}

func init() {
	streamMap = make(map[string]StreamInfo)
}

func GetNowStreaming(c *gin.Context) {

	ChID := c.Param("ChanID")
	temp := streamMap[ChID]
	fmt.Println(streamMap)
	fmt.Println(temp.Url != "" && temp.ExpireTime.Sub(time.Now()).Seconds() > 0)
	if temp.Url != "" && temp.ExpireTime.Sub(time.Now()).Seconds() > 0 {
		c.String(http.StatusOK, temp.Url)
		return
	}

	type cid struct {
		ChannelID string `json:"channelId"`
	}

	filter := expression.Name("customUrl").Equal(
		expression.Value(aws.String(ChID)))
	project := expression.NamesList(expression.Name("channelID"))

	var cc cid

	cl := DB.DbQueryMany(DB.UserTable, filter, project, cc)

	cx := cl[0].(map[string]interface{})

	result := getVideoID("channel", cx["channelID"].(string))

	temp = StreamInfo{result, time.Now()}
	if result != "" {
		temp.ExpireTime = temp.ExpireTime.Add(time.Minute * 3)
	}
	streamMap[ChID] = temp

	c.String(http.StatusOK, temp.Url)
	//c.String(http.StatusOK, "null")

}

func getVideoID(channelType string, channelID string) string {

	var searchUrl string

	if channelType == "user" {
		searchUrl = "https://www.googleapis.com/youtube/v3/search?part=snippet&q=" + channelID + "&type=video&eventType=live&key=" + ytKey
	} else {
		searchUrl = "https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=" + channelID + "&type=video&eventType=live&key=" + ytKey
	}

	res, err := http.Get(searchUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%s\n", sitemap)
	s := string(sitemap)

	var r Reply

	if err := json.Unmarshal([]byte(s), &r); err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Println(r)

	if len(r.Items) > 0 {
		vid := r.Items[0].Id.VideoID
		return vid
	}
	return ""

}

func SetStreamingByUrl(c *gin.Context) {

	chid := c.Query("channel")
	streamUrl := c.Query("videoID")

	if chid != "" && streamUrl != "" {
		streamMap[chid] = StreamInfo{streamUrl, time.Now().Add(time.Minute * 10)}
		c.String(http.StatusOK, "ok")
		return
	}

	c.String(http.StatusOK, "fail")

}

func getVideoIDbyClient(channelType, channelID, clientToken string) string {

	var searchUrl string

	if channelType == "user" {
		searchUrl = "https://www.googleapis.com/youtube/v3/search?part=snippet&q=" + channelID + "&type=video&eventType=live&key=" + clientToken
	} else {
		searchUrl = "https://www.googleapis.com/youtube/v3/search?part=id&channelId=" + channelID + "&type=video&eventType=live&access_token=" + clientToken
	}

	res, err := http.Get(searchUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%s", sitemap)
	s := string(sitemap)

	fmt.Println(s)
	var r Reply

	if err := json.Unmarshal([]byte(s), &r); err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Println(r)
	vid := "null"
	if r.Items != nil {
		vid = r.Items[0].Id.VideoID

	}

	return vid
}
