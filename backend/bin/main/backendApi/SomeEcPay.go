package backendApi

import (
	"Renew/bin/main/misc"
	"crypto/sha256"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

func EcpayTimeString() string {
	tN := time.Now()
	return fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", tN.Year(), tN.Month(), tN.Day(), tN.Hour(), tN.Minute(), tN.Second())
}
func CheckMacCalculation(values url.Values) string {
	var tempString []string
	for key, value := range values {
		tempString = append(tempString, fmt.Sprintf("%s=%s", key, strings.Join(value, "")))
	}
	sort.Strings(tempString)
	sortedTemp := strings.Join(tempString, "&")

	beforeEncode := fmt.Sprintf("%s%s%s", "HashKey="+misc.Config.ECpayHashKey+"&", sortedTemp, "&HashIV="+misc.Config.ECpayHashIV)

	afterEncode := url.Values{"encoded": {beforeEncode}}.Encode()
	afterEncode = strings.ToLower(strings.TrimPrefix(afterEncode, "encoded="))

	sha256seed := sha256.New()
	sha256seed.Write([]byte(afterEncode))

	return strings.ToUpper(fmt.Sprintf("%x", sha256seed.Sum(nil)))
}
