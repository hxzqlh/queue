# queue

## 特征：

* 每个消息在生成时指定消息类型
* 不同类型的消息分发到不同的消息队列
* 不同类型的消息由不同的消息处理器消费
* 消息消费完成后可以执行回调逻辑

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
}, MessageType_MsgA)
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