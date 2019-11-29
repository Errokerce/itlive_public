package backendApi

import (
	DB "Renew/bin/main/aws/database"
	"Renew/bin/main/misc"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	Ep "github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type apiHandler func(c *gin.Context, claims *Claims) interface{}

//上傳 是否為賣家 以及 頻道資訊
func PostIsSeller(c *gin.Context, claims *Claims) interface{} {

	type postBody struct {
		IsSeller  bool   `json:"is_seller"`
		CustomUrl string `json:"custom_url"`
		ChannelID string `json:"channel_id"`
	}

	ps := postBody{}
	_ = c.ShouldBindJSON(&ps)

	fmt.Println(ps)
	//	檢查重複
	if ps.IsSeller {
		filter := Ep.Or(
			Ep.Name("channelID").Equal(
				Ep.Value(aws.String(ps.ChannelID))),
			Ep.Name("customUrl").Equal(
				Ep.Value(aws.String(ps.CustomUrl))))
		project := Ep.NamesList(Ep.Name("channelID"), Ep.Name("customUrl"))

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
	}

	//	投入ＤＢ

	km := map[string]*string{
		"#is": aws.String("isSeller"),
	}

	type updv struct {
		IsSeller bool `json:":is"`
	}

	updcmd := "set #is = :is"

	if ps.IsSeller {
		type updvS struct {
			IsSeller  bool   `json:":is"`
			CustomUrl string `json:":url"`
			ChannelID string `json:":id"`
		}

		km["#url"] = aws.String("customUrl")
		km["#id"] = aws.String("channelID")

		updcmd += " , #url = :url , #id = :id"
		fmt.Println(ps)
		if ps.CustomUrl == "" {
			ps.CustomUrl = claims.UserID
		}
		fmt.Println(ps)
		DB.DbUpdate(DB.UserTable, "userID", claims.UserID, updcmd, km, updvS{ps.IsSeller, ps.CustomUrl, ps.ChannelID})

		_, err := http.Get(fmt.Sprintf("http://localhost:7000/i/%s", ps.CustomUrl))
		if err != nil {
			fmt.Println(err)
		}

	} else {

		DB.DbUpdate(DB.UserTable, "userID", claims.UserID, updcmd, km, updv{ps.IsSeller})
	}

	return "ok"

}

//取得 使用者個資
//@todo Renewed, need Test
func GetProfileEdit(c *gin.Context, claims *Claims) interface{} {

	type respStruct struct {
		UserName  string `json:"userName"`
		NameFirst string `json:"nameFirst"`
		NameLast  string `json:"nameLast"`
		Mail      string `json:"mail"`
		City      string `json:"city"`
		Address   string `json:"address"`
	}

	var resp respStruct

	filter := Ep.Name("userID").Equal(Ep.Value(aws.String(claims.UserID)))

	project := Ep.NamesList(Ep.Name("userName"), Ep.Name("nameFirst"), Ep.Name("nameLast"), Ep.Name("mail"), Ep.Name("city"), Ep.Name("address"))

	var u DB.User

	cc := DB.DbQueryMany(DB.UserTable, filter, project, resp)
	fmt.Println(cc)

	if len(cc) != 0 {
		err := mapstructure.Decode(cc[0], &resp)
		PrintErr(err)
		fmt.Println(resp)
	}

	if debugMode {
		fmt.Println(u)
	}

	return resp

	//DB.DbQuerybyKey(DB.UserTable, "userID", claims.UserID, u)
	//resp = respStruct{u.UserName, u.NameFirst, u.NameLast, u.Mail, u.City, u.Address}

	//return resp

}

//更新 使用者個資
func PostProfileEdit(c *gin.Context, claims *Claims) interface{} {

	type UserData struct {
		NickName  string `json:"nick_name"`
		NameFirst string `json:"name_first"`
		NameLast  string `json:"name_last"`
		Email     string `json:"email"`
		City      string `json:"city"`
		Address   string `json:"address"`
	}

	var ud UserData

	err := c.ShouldBindJSON(&ud)
	if err != nil {
		fmt.Println(err)
	}

	km := map[string]*string{
		"#nn":  aws.String("userName"),
		"#nf":  aws.String("nameFirst"),
		"#nl":  aws.String("nameLast"),
		"#c":   aws.String("city"),
		"#adr": aws.String("address"),
		"#em":  aws.String("mail"),
	}

	type updv struct {
		NickName  string `json:":nn"`
		NameFirst string `json:":nf"`
		NameLast  string `json:":nl"`
		Email     string `json:":em"`
		City      string `json:":c"`
		Address   string `json:":adr"`
	}

	up := updv{
		ud.NickName, ud.NameFirst, ud.NameLast,
		ud.Email, ud.City, ud.Address,
	}

	UpCmd := "set #nn = :nn , #nl = :nl , #nf = :nf , #c = :c , #em = :em , #adr = :adr"

	err2 := DB.DbUpdate(DB.UserTable, "userID", claims.UserID, UpCmd, km, up)

	if err2 {
		return "err"
	}

	jwtToken := FullToken(claims.UserID)
	//jwtToken, _ := GenerateToken(u.UserID, u.UserName)
	c.SetCookie("jwtAccess", jwtToken, cookieTime, "/", cookieDomain, false, true)

	return "ok"

}

//反正就是驗證Token的部分
func VerifyHandler(handler apiHandler) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		var state string
		var resp interface{}

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
			tokenState, claims, err := ParseToken(token)
			if err != nil {
				fmt.Println(err)
			}
			switch tokenState {
			case TokenAccecpt:
				state = "ok"

				resp = handler(c, claims)
			case TokenTimeout:
				state = "timeout"
			case TokenValidationFailed:
				state = "tokenFail"
			}
		}

		c.JSONP(http.StatusOK, gin.H{
			"state": state,
			"data":  resp,
		})

		if debugMode {
			fmt.Println(state)
		}
	}

	return gin.HandlerFunc(fn)
}

//檢查 商品是否都可以購買
//@todo mustRenew
func checkItemAvailable(items []DB.Item) (bool, []DB.Item, string) {
	/*
		filt := expression.Name("available").Equal(expression.Value(aws.Bool(true)))
		proj := expression.NamesList(expression.Name("itemID"), expression.Name("available"))

		ia := database.DbQueryMany(database.ItemTable, filt, proj, database.Item{})
		fmt.Println(ia)*/

	var backI []DB.Item
	var allavailable = true
	var preSeller = ""

	for _, item := range items {
		var temp DB.Item
		DB.DbQuerybyKey(DB.ItemTable, "itemID", item.ItemID, &temp)
		temp.Amount = item.Amount
		backI = append(backI, temp)
		if !temp.Available {
			allavailable = false
		}
		if preSeller == "" {
			preSeller = item.Owner
		} else if preSeller != item.Owner {
			allavailable = false
		}
	}

	return allavailable, backI, preSeller
}

//計算訂單資訊 & 傳給ECPAY
//@todo mustRenew
func sendToEcpay(order DB.Order) string {

	var allNameA []string
	for _, i := range order.Items {
		allNameA = append(allNameA, i.Name)
	}

	payload := url.Values{
		"MerchantID":        {"2000132"},
		"MerchantTradeNo":   {fmt.Sprintf("itLive%s", order.OrderID)},
		"StoreID":           {order.OrderOwner},
		"MerchantTradeDate": {order.Date},
		"PaymentType":       {"aio"},
		"TotalAmount":       {strconv.Itoa(order.PriceSum)},
		"TradeDesc":         {"拍賣商品"},
		"ItemName":          {strings.Join(allNameA, "|")},
		"ReturnURL":         {"http://myDomain.com/api/ec/paymentComfirm"},
		"ChoosePayment":     {"ALL"},
		"EncryptType":       {strconv.Itoa(1)},
	}

	payload.Add("CheckMacValue", CheckMacCalculation(payload))

	resp, err := http.PostForm("https://payment-stage.ecpay.com.tw/SP/CreateTrade", payload)
	if err != nil {
		fmt.Println(err)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	type ecRespone struct {
		MerchantID      string `json:"MerchantID"`
		MerchantTradeNo string `json:"MerchantTradeNo"`
		RtnCode         string `json:"RtnMsg"`
		SPToken         string `json:"SPToken"`
		CheckMacValue   string `json:"CheckMacValue"`
	}
	var ecResp ecRespone
	if err := json.Unmarshal([]byte(contents), &ecResp); err != nil {
		fmt.Println("ERROR:", err)
	}
	fmt.Println("in send to ec")
	fmt.Println(ecResp)
	return ecResp.SPToken

}

//接收 User傳來的訂單資訊 然後呼叫 sendToEcpay 並取回交易Token
//@todo mustRenew, JwtRequire
func OrderConfirm(c *gin.Context) {

	type respItems struct {
		Items []DB.Item `json:"items"`
	}
	var iarna respItems

	c.ShouldBindJSON(&iarna)
	available, items, orderOwner := checkItemAvailable(iarna.Items)
	if available {
		o := DB.NewOrder(items, EcpayTimeString(), "admin", "tw", "ty", "123", orderOwner, "n")
		fmt.Println(o.PriceSum)
		spToken := sendToEcpay(o)
		c.JSON(http.StatusOK, gin.H{
			"state":   "ok",
			"spToken": spToken,
		})

	}
}

//更新WebSocket的ChannelMap用的
func GetChannelList() []string {
	type cid struct {
		ChannelID string `json:"channelId"`
	}

	filter := Ep.AttributeExists(Ep.Name("channelID"))
	project := Ep.NamesList(Ep.Name("customUrl"))

	var c cid

	cl := DB.DbQueryMany(DB.UserTable, filter, project, c)

	var sa []string

	for _, cc := range cl {
		m := cc.(map[string]interface{})
		sa = append(sa, m["customUrl"].(string))
	}

	return sa
	//fmt.Println(cl)
}

//abandon
//@todo 可用於renew 使用者個資
func GetUserAll(userID string) interface{} {
	type UseAble struct {
		UserID     string `json:"userID"`
		UserName   string `json:"userMame"`
		ChannelID  string `json:"channelID"`
		ChannelUrl string `json:"customUrl"`
		IsSeller   bool   `json:"isSeller"`
	}
	filter := Ep.Name("userID").Equal(Ep.Value(aws.String(userID)))
	project := Ep.NamesList(Ep.Name("userID"), Ep.Name("userName"), Ep.Name("isSeller"),
		Ep.Name("channelID"), Ep.Name("customUrl"))

	var backU UseAble

	c := DB.DbQueryMany(DB.UserTable, filter, project, backU)
	fmt.Println(c)

	if len(c) != 0 {
		err := mapstructure.Decode(c[0], &backU)
		PrintErr(err)
		fmt.Println(backU)
	}

	return backU
}

//No respond, Just a function to print Error
func PrintErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//NewTokenCreator
func FullToken(userID string) string {

	type UseAble struct {
		UserID     string `json:"userID"`
		UserName   string `json:"userName"`
		ChannelID  string `json:"channelID"`
		ChannelUrl string `json:"customUrl"`
		IsSeller   bool   `json:"isSeller"`
		Verify     bool   `json:"verify"`
	}
	var ua UseAble

	filter := Ep.Name("userID").Equal(Ep.Value(aws.String(userID)))
	project := Ep.NamesList(Ep.Name("userID"), Ep.Name("userName"), Ep.Name("isSeller"),
		Ep.Name("channelID"), Ep.Name("customUrl"))

	c := DB.DbQueryMany(DB.UserTable, filter, project, UseAble{})

	if len(c) == 0 {
		return "0"
	}

	err := mapstructure.Decode(c[0], &ua)
	PrintErr(err)
	fmt.Println(ua)

	if !ua.IsSeller {
		ua.ChannelID = "0"
		ua.ChannelUrl = "0"
	}

	jwtToken, _ := GenerateToken(ua.UserID, ua.UserName, ua.ChannelID, ua.ChannelUrl, ua.IsSeller)
	return jwtToken
}

//@todo postItem, need Test
func PostItem(c *gin.Context, claims *Claims) interface{} {
	type tempItem struct {
		ItemName     string `json:"itemName"`
		ItemQuantity int    `json:"itemQuantity"`
		ASetQuantity int    `json:"aSetQuantity"`
		ItemPrice    int    `json:"itemPrice"`
		AddPrice1    int    `json:"addPrice1"`
		AddPrice2    int    `json:"addPrice2"`
		AddPrice3    int    `json:"addPrice3"`
		ItemText     string `json:"itemText"`
	}
	var temp tempItem
	km := make(map[string]*string)
	km["#owner"] = aws.String("owner")
	UpCmd := "set #owner = :owner, "

	PrintErr(c.ShouldBindJSON(&temp))
	//keyNames := structs.Names(temp)

	km["#itemName"] = aws.String("name")
	km["#itemQuantity"] = aws.String("amount")
	km["#aSetQuantity"] = aws.String("aset_quantity")
	km["#itemPrice"] = aws.String("price")
	km["#addPrice1"] = aws.String("addPrice_1")
	km["#addPrice2"] = aws.String("addPrice_2")
	km["#addPrice3"] = aws.String("addPrice_3")
	km["#itemText"] = aws.String("describe")
	UpCmd += "#itemName = :itemName, "
	UpCmd += "#itemQuantity = :itemQuantity, "
	UpCmd += "#aSetQuantity = :aSetQuantity, "
	UpCmd += "#itemPrice = :itemPrice, "
	UpCmd += "#addPrice1 = :addPrice1, "
	UpCmd += "#addPrice2 = :addPrice2, "
	UpCmd += "#addPrice3 = :addPrice3, "
	UpCmd += "#itemText = :itemText, "

	UpCmd = strings.TrimSuffix(UpCmd, ", ")

	type updv struct {
		ItemName     string `json:":itemName"`
		ItemQuantity int    `json:":itemQuantity"`
		ASetQuantity int    `json:":aSetQuantity"`
		ItemPrice    int    `json:":itemPrice"`
		AddPrice1    int    `json:":addPrice1"`
		AddPrice2    int    `json:":addPrice2"`
		AddPrice3    int    `json:":addPrice3"`
		ItemText     string `json:":itemText"`
		Owner        string `json:":owner"`
	}

	up := updv{
		temp.ItemName,
		temp.ItemQuantity,
		temp.ASetQuantity,
		temp.ItemPrice,
		temp.AddPrice1,
		temp.AddPrice2,
		temp.AddPrice3,
		temp.ItemText,
		claims.UserID,
	}

	fmt.Printf("temp:\n%v\n", temp)
	fmt.Printf("km:\n%v\n", km)
	fmt.Printf("cmd:\n%v\n", UpCmd)
	fmt.Printf("up:\n%v\n", up)

	thisID := fmt.Sprintf("%08d", misc.GetLatestItemID())

	err2 := DB.DbUpdate(DB.ItemTable, "itemID", thisID, UpCmd, km, up)
	if err2 {
		return "err"
	}
	return gin.H{"itemID": thisID}
}

//testFunc abandon next commit
func TestPostItem(b []byte) {

	type tempItem struct {
		ItemName     string `json:"itemName,omitempty"`
		ItemQuantity int    `json:"itemQuantity,omitempty"`
		ASetQuantity int    `json:"aSetQuantity,omitempty"`
		ItemPrice    int    `json:"itemPrice,omitempty"`
		AddPrice1    int    `json:"addPrice1,omitempty"`
		AddPrice2    int    `json:"addPrice2,omitempty"`
		AddPrice3    int    `json:"addPrice3,omitempty"`
		ItemText     string `json:"itemText,omitempty"`
	}
	var temp tempItem
	var keyName map[string]interface{}
	var km map[string]*string
	UpCmd := "set "

	for k, _ := range keyName {
		switch k {
		case "itemName":
			km["#itemName"] = aws.String("name")
			UpCmd += "#itemName = :itemName, "
		case "itemQuantity":
			km["#itemQuantity"] = aws.String("amount")
			UpCmd += "#itemQuantity = :itemQuantity, "
		case "aSetQuantity":
			km["#aSetQuantity"] = aws.String("aset_quantity")
			UpCmd += "#aSetQuantity = :aSetQuantity, "
		case "itemPrice":
			km["#itemPrice"] = aws.String("price")
			UpCmd += "#itemPrice = :itemPrice, "
		case "addPrice1":
			km["#addPrice1"] = aws.String("addPrice_1")
			UpCmd += "#addPrice1 = :addPrice1, "
		case "addPrice2":
			km["#addPrice2"] = aws.String("addPrice_2")
			UpCmd += "#addPrice2 = :addPrice2, "
		case "addPrice3":
			km["#addPrice3"] = aws.String("addPrice_3")
			UpCmd += "#addPrice3 = :addPrice3, "
		case "itemText":
			km["#itemText"] = aws.String("describe")
			UpCmd += "#itemText = :itemText, "
		}
	}

	UpCmd = strings.TrimSuffix(UpCmd, ", ")

	type updv struct {
		ItemName     string `json:":itemName,omitempty"`
		ItemQuantity int    `json:":itemQuantity,omitempty"`
		ASetQuantity int    `json:":aSetQuantity,omitempty"`
		ItemPrice    int    `json:":itemPrice,omitempty"`
		AddPrice1    int    `json:":addPrice1,omitempty"`
		AddPrice2    int    `json:":addPrice2,omitempty"`
		AddPrice3    int    `json:":addPrice3,omitempty"`
		ItemText     string `json:":itemText,omitempty"`
	}

	up := updv{
		temp.ItemName,
		temp.ItemQuantity,
		temp.ASetQuantity,
		temp.ItemPrice,
		temp.AddPrice1,
		temp.AddPrice2,
		temp.AddPrice3,
		temp.ItemText,
	}

	fmt.Printf("\n%v\n", UpCmd)
	fmt.Printf("\n%v\n", km)
	fmt.Printf("\n%v\n", up)

	err2 := DB.DbUpdate(DB.ItemTable, "itemID", strconv.Itoa(misc.GetLatestItemID()), UpCmd, km, up)
	if err2 {
		return
	}

}

//成交後確認最高價＆得標者
//@todo need test
func ConfirmHighestPrice(c *gin.Context, claims *Claims) interface{} {

	type tempS struct {
		ItemID       string `json:"item_id"`
		HighestPrice int    `json:"highest_price"`
		//Winner       string `json:"winner"`
	}
	var temp tempS
	PrintErr(c.ShouldBindJSON(&temp))

	UpCmd := "set #highest_price=:highest_price" //, #winner=:winner

	km := map[string]*string{
		"#highest_price": aws.String("highest_price"),
		//"#winner":        aws.String("winner"),
	}

	type updv struct {
		HighestPrice int `json:":highest_price"`
		//Winner       string `json:":winner"`
	}

	err2 := DB.DbUpdate(DB.ItemTable, "itemID", temp.ItemID, UpCmd, km, updv{temp.HighestPrice})
	if err2 {
		return "err"
	}
	return "ok"

}

func ConfirmOrder(c *gin.Context, claims *Claims) interface{} {
	type (
		shipping struct {
			Address     string `json:"address"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			Country     string `json:"country"`
		}
		oItem struct {
			HeightPrice string `json:"heightprice"`
			Id          string `json:"id"`
			ItemName    string `json:"itemName"`
		}
		RecvOrder struct {
			ShippingInfo shipping `json:"shipping_info"`
			Items        []oItem  `json:"items"`
		}
	)

	var o RecvOrder

	PrintErr(c.ShouldBindJSON(&o))

	ia := make([]DB.Item, 0)
	for _, v := range o.Items {
		i := DB.Item{}
		i.ItemID = v.Id
		i.Name = v.ItemName
		price, err := strconv.Atoi(v.HeightPrice)
		PrintErr(err)
		i.Price = price
		i.Amount = 1

		ia = append(ia, i)

	}

	oo := DB.NewOrder(ia, EcpayTimeString(), o.ShippingInfo.Name, "tw", o.ShippingInfo.Country, o.ShippingInfo.Address, "102655569137082735940", claims.UserID)

	DB.DbAdd(DB.OrderTable, oo)

	spToken := sendToEcpay(oo)
	if spToken != "" {
		return gin.H{
			"state":   "ok",
			"spToken": spToken,
		}
	}
	return gin.H{
		"state": "fail",
	}
}

func TestPayment(c *gin.Context) {
	Items := make([]DB.Item, 10)
	i := DB.Item{}
	thisID := fmt.Sprintf("%08d", misc.GetLatestItemID())
	i.ItemID = thisID

	i.Name = "Pen"
	i.Amount = 1
	i.Price = 500
	Items = append(Items, i)

	o := DB.NewOrder(Items, EcpayTimeString(), "admin", "tw", "桃園市", "平鎮區", "102655569137082735940", "102655569137082735940")

	spToken := sendToEcpay(o)
	c.JSON(http.StatusOK, gin.H{
		"state":   "ok",
		"spToken": spToken,
	})
}
