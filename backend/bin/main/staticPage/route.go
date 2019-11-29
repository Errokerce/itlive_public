package staticPage

import (
	api "Renew/bin/main/backendApi"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"Renew/bin/main/chatroom"
	ws "Renew/bin/main/wsHandler"
)

var r *gin.Engine

func init() {
	rr := gin.Default()
	r = rr

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://www.itlive.nctu.me", "https://websocket.itlive.nctu.me", "https://api.itlive.nctu.me"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"X-Total-Count"}

	r.Use(cors.New(config))

}



func backEndAPI() {

	//r.GET("/api/jwt", api.VerifyJWT)
	//r.GET("/api/self", api.ViewUserdata)

	//r.POST("/login/verify/phone", api.PostPhone)
	//r.POST("/login/verify/code", api.PostCode)

	r.POST("/login/verify/phone", api.VerifyHandler(api.PostPhoneV2))
	r.POST("/login/verify/code", api.VerifyHandler(api.PostCodeV2))

	r.POST("/login/register/isSeller", api.VerifyHandler(api.PostIsSeller))

	r.GET("/login/profileEdit", api.VerifyHandler(api.GetProfileEdit))
	r.POST("/login/profileEdit", api.VerifyHandler(api.PostProfileEdit))

	r.GET("/login/GoogleLogin", api.HandleGoogleLogin)
	r.GET("/login/GoogleCallback", api.HandleGoogleCallback)

	r.GET("/login/GetStreamSrc/:ChanID", api.GetNowStreaming)

	r.GET("/login/updateCache", api.SetStreamingByUrl)

	r.POST("/login/Item", api.VerifyHandler(api.PostItem))
	r.POST("/login/ConfirmHigh", api.VerifyHandler(api.ConfirmHighestPrice))

	r.GET("/login/FakePayment", api.TestPayment)

	r.POST("/login/ConfirmOrder", api.VerifyHandler(api.ConfirmOrder))
	//r.POST("/api/OrderConfirm", api.OrderConfirm)
	//r.GET("/api/newWay", api.VerifyHandler(api.DDD))

}

func hostSocket() {
	r.GET("/s/:ChanID", ws.SocketVer3(chatroom.Rmap))

	r.GET("/i/:customUrl", ws.AddChannel(chatroom.Rmap))
}

func Route(port string) {

	switch port {
	case ":8010":
		backEndAPI()
	case ":7000":
		hostSocket()
	}

	log.Fatal(r.Run(port))
}
