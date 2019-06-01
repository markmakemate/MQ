package MQ

import (
	"MQ/Cache"
	"MQ/Consumer"
	"MQ/Producer"
	"MQ/Utils"
	"sync"
)

/*高并发消息队列的基本想法：
* 并发tcp请求将请求线程先放在一个线程池里，Puber 和Suber分别管理两个线程池。请求是异步并发执行的，
Producer和Consumer端都不需要等待服务器返回信息。Producer生产消息，每个Puber发布一个Block消息是一个原子操作，
每个Suber订阅一个Block的消息是一个原子操作，当队列满时，Puber的goroutine阻塞，当队列空时，Suber的goroutine阻塞。线程池内的Worker
可以手动结束请求。也可以在消息队列关闭后自动GC回收。当消息队列关闭或程序发出关闭信号，则将Puber/Suber关闭。
*/
type MQ struct{
	Puber Producer.Publisher
	Suber Consumer.Subscriber
}
func (q *MQ) Start(topic string, addrOfProd string, addrOfConsume string, size int64, PuberSignal chan bool,
	SuberSignal chan bool){
	//全局读写锁，保证Worker消息操作的原子性
	var mutex = new(sync.RWMutex)
	q.Suber.New(0, topic)
	q.Puber.New(0, topic)
	ProducerListener := Utils.Listen(addrOfProd)
	ConsumerListener := Utils.Listen(addrOfConsume)
	queue := make(chan Cache.Block, size)
	go q.Puber.Start(ProducerListener, queue, PuberSignal, mutex)
	go q.Suber.Start(ConsumerListener, queue, SuberSignal, mutex)
}