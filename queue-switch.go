package queue

import "errors"

const (
	MessageType_MsgCount MessageType = 64 // 最多支持消息类型总数
)

var ErrMessageType = errors.New("ErrMessageType")
var ErrMessageQueue = errors.New("ErrMessageQueue")
var ErrMessageHandle = errors.New("ErrMessageHandle")

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

// 消息发送出去后将被路由到 msgType 对应的队列，消息处理完后执行回调函数 fn
func (s *QueueSwitch) NewMessage(data interface{}, fn MessageCB, msgType MessageType) (*Message, error) {
	if msgType >= MessageType_MsgCount {
		return nil, ErrMessageType
	}
	if s.handles[msgType] == nil {
		return nil, ErrMessageHandle
	}
	if s.queues[msgType] == nil {
		return nil, ErrMessageQueue
	}

	var msg Message
	msg.Type = msgType
	msg.Fn = s.createCB(msgType, fn)
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

func (s *QueueSwitch) createCB(msgType MessageType, fn MessageCB) func(*Message) {
	return func(msg *Message) {
		// 执行自己的 handle
		msg, err := s.handles[msgType](msg)
		if fn != nil {
			// 回调处理结果
			fn(msg, err)
		}
	}
}
