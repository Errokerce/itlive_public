package chatroom

import (
	"time"
)

const (
	//訊息類型 依序對應從0開始的整數
	EVENT_TYPE_MSG    = iota //0
	EVENT_TYPE_JOIN          //1
	EVENT_TYPE_TYPING        //2 etc.
	EVENT_TYPE_LEAVE
	EVENT_TYPE_IMAGE
)

type Event struct {
	Type      int    `json:"type"`
	User      string `json:"user"`
	Timestamp int64  `json:"timestamp"`
	Text      string `json:"text"`
}

func newEvent(typ int, user, msg string) Event {
	return Event{typ, user, time.Now().UnixNano() / 1e6, msg}
}

type Subscription struct {
	id int64

	username string

	Pipe <-chan Event

	emit chan Event

	leave chan int64
}

func (s *Subscription) Leave() {
	s.leave <- s.id
}

func (s *Subscription) Say(message string) {
	s.emit <- newEvent(EVENT_TYPE_MSG, s.username, message)
}
