package database

import (
	"Renew/bin/main/misc"
	"fmt"
	"time"
)

type User struct {
	UserID      string  `json:"userID"`
	UserName    string  `json:"userName,omitempty"`
	Verified    bool    `json:"verify"`
	IsSeller    bool    `json:"isSeller"`
	NameFirst   string  `json:"nameFirst,omitempty"`
	NameLast    string  `json:"nameLast,omitempty"`
	Phone       int     `json:"phone,omitempty"`
	Mail        string  `json:"mail,omitempty"`
	Country     string  `json:"country,omitempty"`
	City        string  `json:"city,omitempty"`
	Address     string  `json:"address,omitempty"`
	Rank        int     `json:"rank"`
	Purchased   []Order `json:"purchased,omitempty"`
	SellingItem []Item  `json:"sellingItem,omitempty"`
	Orders      []Order `json:"orders,omitempty"`
	ChannelID   string  `json:"channelID,omitempty"`
}

type Order struct {
	OrderID       string `json:"orderID"`
	IsPaid        bool   `json:"isPaid"`
	PriceSum      int    `json:"priceSum"`
	Date          string `json:"date"`
	Name          string `json:"name"`
	Country       string `json:"country"`
	City          string `json:"city"`
	Address       string `json:"address"`
	Items         []Item `json:"items"`
	OrderOwner    string `json:"order_owner"`
	Enable        bool   `json:"enable"`
	Buyer         string
	Statue        string
	PaymentStatue string
	ShipStatus    string
	estbTime      time.Time
}

type Item struct {
	ItemID       string `json:"itemID"`
	Name         string `json:"name"`
	Describe     string `json:"describe,omitempty"`
	Price        int    `json:"price"`
	Amount       int    `json:"amount,omitempty"`
	Stock        int    `json:"stock,omitempty"`
	Sold         int    `json:"sold,omitempty"`
	Owner        string `json:"owner"`
	Available    bool   `json:"available"`
	Enable       bool   `json:"enable"`
	AddPrice1    int    `json:"addPrice_1"`
	AddPrice2    int    `json:"addPrice_2"`
	AddPrice3    int    `json:"addPrice_3"`
	asetQuantity int    `json:"aset_quantity"`
}

func NewOrder(items []Item, date, name, country, city, address, owner, buyer string) Order {
	o := Order{}

	o.OrderID = fmt.Sprintf("%08d", misc.GetLatestOrderID())
	o.IsPaid = false
	o.PriceSum = 0
	for _, i := range items {
		o.PriceSum += i.Price * i.Amount
	}
	o.Date = date
	o.Name = name
	o.Country = country
	o.City = city
	o.Address = address
	o.Items = items
	o.OrderOwner = owner
	o.Buyer = buyer

	return o
}

func NewUser(id string, fn, ln, em string) User {

	u := User{
		id,
		fn + ln,
		false,
		false,
		fn,
		ln,
		0,
		em,
		"Tw",
		"",
		"",
		5,
		nil,
		nil,
		nil,
		"",
	}
	return u
}

func NewItem(owner, name string, price, stock int, available bool) Item {

	i := Item{}
	i.ItemID = fmt.Sprintf("%08d", misc.GetLatestItemID())
	i.Name = name
	i.Price = price
	i.Stock = stock
	i.Available = available
	i.Owner = owner

	return i
}
