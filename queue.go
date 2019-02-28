package queue

type MessageType int32

type Message struct {
	Data interface{}        // 消息内容
	Fn   func(msg *Message) // 消息处理函数
	Type MessageType        // 消息类型
}

type IQueue interface {
	Send(msg *Message)
}

type Queue struct {
	ch chan *Message
}

func (q *Queue) Send(msg *Message) {
	q.ch <- msg
}

func NewQueue() *Queue {
	q := &Queue{make(chan *Message, 2048)}
	go func() {
		var msgs [1024]*Message
		for {
			msgs[0] = <-q.ch

			length := len(q.ch) + 1
			if length > 1024 {
				length = 1024
			}
			for i := 1; i < length; i++ {
				msgs[i] = <-q.ch
			}
			for i := 0; i < length; i++ {
				msgs[i].Fn(msgs[i])
			}
		}
	}()
	return q
}
