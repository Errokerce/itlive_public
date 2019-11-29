package main

import (
	api "Renew/bin/main/backendApi"
	"Renew/bin/main/staticPage"
	. "fmt"
	"net/url"
	"strconv"
)

const (
	WebPort    = 8080
	VERSION    = 1
	LoginPort  = ":8010"
	SocketPort = ":7000"
)

func main() {

	// Host Server的部分
	staticPage.Route(LoginPort)
	//staticPage.Route(SocketPort)


}
