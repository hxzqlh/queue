package queue

type MessageType int32

type Message struct {
	Data interface{}        // 消息内容
	Fn   func(msg *Message) // 消息处理函数
	Type MessageType        // 消息类型
}

type IQueue interface {
	Send(msg *Message)
	GetLevel() int
}

type Queue struct {
	ch    chan *Message
	level int // 队列优先级，值越小，优先级越高
}

func (q *Queue) Send(msg *Message) {
	q.ch <- msg
}

func (q *Queue) GetLevel() int {
	return q.level
}

func NewQueue(level int) *Queue {
	q := &Queue{make(chan *Message, 2048), level}
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
