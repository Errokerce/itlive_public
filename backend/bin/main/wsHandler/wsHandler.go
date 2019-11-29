package wsHandler

import (
	api "Renew/bin/main/backendApi"
	"Renew/bin/main/chatroom"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ROOM = chatroom.RoomS

func AddChannel(chMap map[string]*chatroom.Room) gin.HandlerFunc {
	return func(c *gin.Context) {
		ChID := c.Param("ChanID")

		temp := chatroom.NewRoom()
		chMap[ChID] = temp
		c.Status(http.StatusOK)
	}

}

func SocketVer3(chMap map[string]*chatroom.Room) gin.HandlerFunc {
	return func(c *gin.Context) {

		//取得網址中的頻道名稱
		ChID := c.Param("ChanID")

		//設定一個頻道變數 但是這裡不能放空值的，所以把主頻道放進來
		var ROOM = chatroom.RoomS

		//如果有抓到頻道ＩＤ就在ＭＡＰ裡查查看
		if ChID != "" {
			var EXIST bool
			ROOM, EXIST = chMap[ChID]
			if !EXIST { // 如果不存在的話再把主頻道丟回去ＲＯＯＭ, 但通常不會到這一步
				ROOM = chatroom.RoomS
			}
		}

		var name string
		token, err := c.Cookie("jwtAccess")
		if err != nil {
			fmt.Println(err)
		}

		if token == "" {
			name = "notLogin"
		} else {
			tokenState, claims, err := api.ParseToken(token)
			if err != nil {
				fmt.Println(err)
			}
			switch tokenState {
			case api.TokenAccecpt:
				name = claims.UserName
			}
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}

		// 加入房间
		evs := ROOM.GetArchive()
		if name != "notLogin" {
			ROOM.MsgJoin(name)
		}
		control := ROOM.Join(name)
		defer control.Leave()

		// 先把历史消息推送出去
		for _, event := range evs {
			if conn.WriteJSON(&event) != nil {
				// 用户断开连接
				return
			}
		}

		// 开启通道监听用户事件然后发送给聊天室
		newMessages := make(chan string)
		if name != "notLogin" {
			go func() {
				var res = struct {
					Msg string `json:"msg"`
				}{}
				for {
					err := conn.ReadJSON(&res)
					if err != nil {
						// 用户断开连接
						close(newMessages)
						return
					}
					newMessages <- res.Msg
				}
			}()
		}

		// 接收消息，在这里阻塞请求，循环退出就表示用户已经断开
		for {
			select {
			case event := <-control.Pipe:
				if conn.WriteJSON(&event) != nil {
					// 用户断开连接
					return
				}
			case msg, ok := <-newMessages:
				// 断开连接
				if !ok {
					return
				}
				control.Say(msg)
			}
		}
	}
}
