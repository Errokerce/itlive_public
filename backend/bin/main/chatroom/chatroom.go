package chatroom

import (
	api "Renew/bin/main/backendApi"
	"container/list"
	"fmt"
	"math/rand"
	"time"
)

const archiveSize = 20
const chanSize = 10
const MSG_JOIN = `{"msg":"加入房間","type":"ConsoleMsg"}`
const MSG_LEAVE = `{"msg":"離開房間","type":"ConsoleMsg"}`
const MSG_TYPING = "[正在输入]"

var RoomS *Room

var Rmap map[string]*Room

// 聊天室
type Room struct {
	users     map[int64]chan Event
	userName  map[int64]string
	userCount int64
	idx       int64

	publishChn chan Event

	archive *list.List

	archiveChan chan chan []Event

	joinChn chan chan Subscription

	leaveChn chan int64
}

func NewRoom() *Room {
	r := &Room{
		users:     map[int64]chan Event{},
		userName:  map[int64]string{}, //@@@@@@@@@@@@@@@@@@@@@@@@@@
		userCount: 0,
		idx:       0,

		publishChn:  make(chan Event, chanSize),
		archiveChan: make(chan chan []Event, chanSize),
		archive:     list.New(),

		joinChn:  make(chan chan Subscription, chanSize),
		leaveChn: make(chan int64, chanSize),
	}

	go r.Serve()

	return r
}

func (r *Room) MsgJoin(user string) {
	r.publishChn <- newEvent(EVENT_TYPE_JOIN, user, MSG_JOIN)
}

func (r *Room) MsgSay(user, message string) {
	r.publishChn <- newEvent(EVENT_TYPE_MSG, user, message)
}

func (r *Room) MsgLeave(user string) {
	r.publishChn <- newEvent(EVENT_TYPE_MSG, user, MSG_LEAVE)
}

func (r *Room) Remove(id int64) {
	r.leaveChn <- id
}

func (r *Room) Join(username string) Subscription {
	resp := make(chan Subscription)
	r.joinChn <- resp
	s := <-resp
	s.username = username

	r.userName[r.idx] = username

	r.MsgSay("SYSTEM", fmt.Sprintf(`{"online":%d}`, r.userCount))

	return s
}

func (r *Room) GetArchive() []Event {
	ch := make(chan []Event)
	r.archiveChan <- ch
	return <-ch
}

func (r *Room) Serve() {
	for {
		select {

		case ch := <-r.joinChn:
			chE := make(chan Event, chanSize)
			r.userCount++
			r.idx++
			r.users[r.idx] = chE
			ch <- Subscription{
				id:    r.idx,
				Pipe:  chE,
				emit:  r.publishChn,
				leave: r.leaveChn,
			}
		case arch := <-r.archiveChan:
			events := []Event{}
			for e := r.archive.Front(); e != nil; e = e.Next() {
				events = append(events, e.Value.(Event))
			}
			arch <- events
		case event := <-r.publishChn:
			for _, v := range r.users {
				v <- event

			}
			if r.archive.Len() >= archiveSize {
				r.archive.Remove(r.archive.Front())
			}
			r.archive.PushBack(event)
		case k := <-r.leaveChn:
			if _, ok := r.users[k]; ok {
				r.MsgLeave(r.userName[k]) //新增這個		@@@@@@@@@@@@@@@@@@@@@@@@@@
				delete(r.users, k)
				delete(r.userName, k) //還有這個			@@@@@@@@@@@@@@@@@@@@@@@@@@
				r.userCount--
			}
		}
	}
}

//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
func (r *Room) RndUser() string {
	ul := make(map[int64]string)
	count := int64(0)
	rand.Seed(time.Now().UnixNano())
	for _, u := range r.userName {
		ul[count] = u
		count++
	}
	return ul[rand.Int63n(count)]
}

func (r *Room) NameList() []string {

	var nl []string
	for _, u := range r.userName {
		fmt.Println("add")
		nl = append(nl, u)
	}
	return nl
}

//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

// 开启goroutine loop Serve
func init() {
	Rmap = make(map[string]*Room)

	RoomS = NewRoom()
	cl := api.GetChannelList()

	for _, c := range cl {
		temp := NewRoom()
		Rmap[c] = temp
	}

}
