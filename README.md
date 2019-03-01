# queue-switch 2.0

## 1.0 特征：

* 每个消息在生成时指定消息类型
* 不同类型的消息分发到不同的消息队列
* 不同类型的消息由不同的消息处理器消费
* 消息消费完成后可以执行回调逻辑

2.0 主要新增了一个功能：

* 支持多优先级消息队列

在某些业务系统中，接收到消息后，为了确保系统宕机后，之前未处理成功的消息也能恢复出来，一般是先把消息写入一个日志队列，成功后，然后再把消息写入业务逻辑队列。

来看下具体的代码改动：

`IQueue` 接口增加返回优先级方法

```
type IQueue interface {
	Send(msg *Message)
	GetLevel() int
}
```

规定：`GetLevel()` 返回的数值越小，优先级越高。

消息路由器生成新消息时，可以指定该消息路由的顺序：从优先级高的队列链式执行到优先级低的队列，处理完成后执行回调函数 cb。

这里，`msgTypes` 必须保持增序，比如，`[]{1,2,3}`， 同优先级的可以并列出现，比如，`[]{1,2,2,3}`

```
func (s *QueueSwitch) NewMessage(data interface{}, cb MessageCB, msgTypes ...MessageType) (*Message, error)
```

## 使用说明

初始化流程：

* `qs := NewQueueSwitch()`生成消息路由器实例
* 调用 `qs.SetRoute`和`qs.SetHandle` 对不同类型的消息设置消息路由和消息处理器

消息发送流程：

* 构造消息：

```
ch := make(chan interface{}, 1) // 消息处理结果会写入该管道
msg, _ = qw.NewMessage("Hello World!", func(msg *Message, err error) {
	ch <- msg.Data
}, MessageType_MsgA, MessageType_MsgB)
```

* 发送消息到路由器：

```
qw.Send(msg)
```

* 等待消息处理完成：

```
result := <-ch
fmt.Println(result)
```
