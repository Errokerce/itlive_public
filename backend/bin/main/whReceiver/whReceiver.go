package whReceiver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

func main() {
	r := gin.Default()

	r.GET("/channelEvent", func(c *gin.Context) {
		challenge, _ := c.GetQuery("hub.challenge")
		c.String(200, challenge)
	})
	r.POST("/channelEvent", func(c *gin.Context) {

		byteBody, _ := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()

		fmt.Println("Hook string ::: ", string(byteBody))

		c.String(200, "")
	})

	log.Fatal(r.Run(":8001"))
}
