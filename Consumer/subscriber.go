package Consumer

import (
	"MQ/Cache"
	"MQ/Utils"
	"MQ/Worker"
	"net"
	"os"
	"sync"
)
//订阅者有独立的id和topic
type Subscriber struct {
	topic string
	id int
	//全局读写锁，保证Worker Produce消息操作的原子性
	RwMutex *sync.RWMutex
}
func (s *Subscriber) New(id int, topic string, mutex *sync.RWMutex){
	s.id = id
	s.topic = topic
	s.RwMutex = mutex
}
func (s *Subscriber) GetTopic() string{
	return s.topic
}
func (s *Subscriber) SetTopic(topic string){
	s.topic = topic
}
func (s *Subscriber) GetId() int{
	return s.id
}
func (s *Subscriber) SetId(id int){
	s.id = id
}

func (s *Subscriber) Sub(worker Worker.AbstractWorker, queue chan Cache.Block){
	//队列空 阻塞
	v := <- queue
    //读写锁，保证Worker consume消息的操作原子性
	s.RwMutex.RLock()
	worker.Consume(v)
	s.RwMutex.RUnlock()
}
func (s *Subscriber) Start(listener net.Listener, queue chan Cache.Block, sign chan bool){
	var wpool = new(Worker.Workerpool)
	var rwMutex = new(sync.RWMutex)
	defer listener.Close()
	wpool.Init(1000000)
	//通道关闭则退出Suber，收到关闭信号量则退出Suber。若出现未知错误，通道均阻塞，则退出Suber
	select {
	case _, ok :=<- queue:
		if !ok{
			return
		}else{
			for{
				conn, err := listener.Accept()
				Utils.ChkError(err)
				go Utils.Handler(conn, wpool)
				for x:=range wpool.GetChannel(){
					worker := x
					go Sub(worker, queue, rwMutex)
				}
			}
		}
	case s :=<- sign:
		if s{
			os.Exit(0)
		}
	default:
		os.Exit(3)
	}


}
