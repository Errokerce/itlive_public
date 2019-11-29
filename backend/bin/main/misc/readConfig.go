package misc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var filePath = "config/config.json"
var Config ConfigFile

type ConfigFile struct {
	UserID        int    `json:"user_latest_id"`
	ItemLatestID  int    `json:"item_latest_id"`
	OrderLatestID int    `json:"order_latest_id"`
	SmsAcc        string `json:"sms_account"`
	SmsPwd        string `json:"sms_password"`
	JwtSecret     string `json:"jwt_secret"`
	PhoneSecret   string `json:"phone_secret"`
	ECpayHashKey  string `json:"ecpay_hash_key"`
	ECpayHashIV   string `json:"ecpay_hash_iv"`
}

func init() {
	readIDfile()
}

func readIDfile() {

	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &Config)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func writeIDfile() {

	file, _ := json.MarshalIndent(&Config, "", " ")

	_ = ioutil.WriteFile(filePath, file, 0644)

}

func GetLatestItemID() int {
	Config.ItemLatestID += 1
	writeIDfile()
	return Config.ItemLatestID
}
func GetLatestOrderID() int {
	Config.OrderLatestID += 1
	writeIDfile()
	return Config.OrderLatestID
}
