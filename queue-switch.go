package queue

import "errors"

const (
	MessageType_MsgCount MessageType = 64 // 最多支持消息类型总数
)

var ErrMessageType = errors.New("ErrMessageType")
var ErrMessageQueue = errors.New("ErrMessageQueue")
var ErrMessageHandle = errors.New("ErrMessageHandle")
var ErrMessageLevel = errors.New("ErrMessageLevel")
var ErrEmptyType = errors.New("ErrEmptyType")

type MessageHandle func(*Message) (*Message, error)
type MessageCB func(*Message, error)

type QueueSwitch struct {
	handles []MessageHandle
	queues  []IQueue
}

func NewQueueSwitch() *QueueSwitch {
	return &QueueSwitch{
		handles: make([]MessageHandle, MessageType_MsgCount),
		queues:  make([]IQueue, MessageType_MsgCount),
	}
}

// 消息发送出去后将被路由到 msgType 对应的队列
// 按照队列优先级从高到低的顺序链式执行到最后的消息队列
// 消息处理完后执行回调函数 cb
func (s *QueueSwitch) NewMessage(data interface{}, cb MessageCB, msgTypes ...MessageType) (*Message, error) {
	if len(msgTypes) == 0 {
		return nil, ErrEmptyType
	}
	for _, msgType := range msgTypes {
		if msgType >= MessageType_MsgCount {
			return nil, ErrMessageType
		}
		if s.handles[msgType] == nil {
			return nil, ErrMessageHandle
		}
		if s.queues[msgType] == nil {
			return nil, ErrMessageQueue
		}
	}
	var fn func(*Message)
	for i := len(msgTypes) - 1; i >= 0; i-- {
		if i == len(msgTypes)-1 {
			fn = s.createLastFn(msgTypes[i], cb)
		} else {
			level1 := s.queues[msgTypes[i]].GetLevel()
			level2 := s.queues[msgTypes[i+1]].GetLevel()
			if level1 > level2 {
				return nil, ErrMessageLevel
			}
			fn = s.createFn(msgTypes[i], msgTypes[i+1], fn, cb)
		}
	}

	var msg Message
	msg.Type = msgTypes[0]
	msg.Fn = fn
	msg.Data = data
	return &msg, nil
}

// 设置消息路由，根据消息类型分发到对应类型的队列中
func (s *QueueSwitch) SetRoute(msgType MessageType, q IQueue) {
	if msgType >= MessageType_MsgCount {
		panic("err msg type")
	}
	s.queues[int(msgType)] = q
}

// 绑定消息处理器，不同类型的消息由不同的处理器处理
func (s *QueueSwitch) SetHandle(msgType MessageType, fn MessageHandle) {
	if msgType >= MessageType_MsgCount {
		panic("err msg type")
	}
	s.handles[int(msgType)] = fn
}

func (s *QueueSwitch) Send(msg *Message) error {
	if msg.Type >= MessageType_MsgCount {
		return ErrMessageType
	}
	q := s.queues[int(msg.Type)]
	if q == nil {
		return ErrMessageQueue
	}
	q.Send(msg)
	return nil
}

func (s *QueueSwitch) createFn(msgType, nextType MessageType, nextFn func(*Message), cb MessageCB) func(*Message) {
	return func(msg *Message) {
		//执行自己的handle
		msg, err := s.handles[msgType](msg)
		if err != nil {
			if cb != nil {
				cb(msg, err)
			}
			return
		}
		//把消息发送给下一个channel
		msg.Type = nextType
		msg.Fn = nextFn
		s.Send(msg)
	}
}

func (s *QueueSwitch) createLastFn(msgType MessageType, cb MessageCB) func(*Message) {
	return func(msg *Message) {
		// 执行自己的 handle
		msg, err := s.handles[msgType](msg)
		if cb != nil {
			// 回调处理结果
			cb(msg, err)
		}
	}
}
