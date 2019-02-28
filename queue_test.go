package queue

import (
	"fmt"
	"testing"
)

var (
	MessageType_MsgA MessageType = 0
	MessageType_MsgB MessageType = 1
)

type TradeServer struct {
	sw *QueueSwitch
	q  *Queue
}

func NewTradeServer() *TradeServer {
	s := &TradeServer{}
	s.sw = NewQueueSwitch()
	s.q = NewQueue()
	s.RegistMsgA()
	s.RegistMsgB()
	return s
}

func (s *TradeServer) RegistMsgA() {
	s.sw.SetRoute(MessageType_MsgA, s.q)
	s.sw.SetHandle(MessageType_MsgA, func(msg *Message) (*Message, error) {
		fmt.Println("A type msg:", msg)
		// some logic...
		msg.Data = "Hello, MsgA"
		return msg, nil
	})
}

func (s *TradeServer) RegistMsgB() {
	s.sw.SetRoute(MessageType_MsgB, s.q)
	s.sw.SetHandle(MessageType_MsgB, func(msg *Message) (*Message, error) {
		fmt.Println("B type msg:", msg)
		// some logic...
		msg.Data = "Hello, MsgB"
		return msg, nil
	})
}

func TestQueue(t *testing.T) {
	ch := make(chan interface{}, 1)
	s := NewTradeServer()
	var msg *Message

	msg, _ = s.sw.NewMessage("test-A", func(msg *Message, err error) {
		fmt.Println("in msg-A callback")
		ch <- msg.Data
	}, MessageType_MsgA)
	s.sw.Send(msg)
	fmt.Println(<-ch)

	msg, _ = s.sw.NewMessage("test-B", func(msg *Message, err error) {
		fmt.Println("in msg-B callback")
		ch <- msg.Data
	}, MessageType_MsgB)
	s.sw.Send(msg)
	fmt.Println(<-ch)
}
