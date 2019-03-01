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
	sw     *QueueSwitch
	logq   *Queue
	matchq *Queue
}

func NewTradeServer() *TradeServer {
	s := &TradeServer{}
	s.sw = NewQueueSwitch()
	s.logq = NewQueue(1)
	s.matchq = NewQueue(2)
	s.RegistMsgA()
	s.RegistMsgB()
	return s
}

func (s *TradeServer) RegistMsgA() {
	s.sw.SetRoute(MessageType_MsgA, s.logq)
	s.sw.SetHandle(MessageType_MsgA, func(msg *Message) (*Message, error) {
		fmt.Println("recv log msg:", msg)
		return msg, nil
	})
}

func (s *TradeServer) RegistMsgB() {
	s.sw.SetRoute(MessageType_MsgB, s.matchq)
	s.sw.SetHandle(MessageType_MsgB, func(msg *Message) (*Message, error) {
		fmt.Println("recv match msg:", msg)
		// some logic...
		msg.Data = "tx-result"
		return msg, nil
	})
}

func TestQueueSwitch(t *testing.T) {
	ch := make(chan interface{}, 1)
	s := NewTradeServer()

	var msg *Message
	msg, _ = s.sw.NewMessage("tx-data", func(msg *Message, err error) {
		fmt.Println("in callback")
		ch <- msg.Data
	}, MessageType_MsgA, MessageType_MsgB)
	s.sw.Send(msg)
	fmt.Println(<-ch)
}
