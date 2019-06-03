package Producer

import (
	"MQ/Cache"
	"MQ/Utils"
	"MQ/Worker"
	"net"
	"os"
	"sync"
)

//发布者有独立的id和发布的topic
type Publisher struct {
	id int
	topic string
	//全局读写锁，保证Worker Produce消息操作的原子性

}
func (p *Publisher) New(id int, topic string){
	p.id = id
	p.topic = topic
}
func (p *Publisher) GetId() int{
	return p.id
}
func (p *Publisher) SetId(id int) {
	p.id = id
}
func (p *Publisher) SetTopic(topic string){
	p.topic = topic
}
func (p *Publisher) GetTopic() string{
	return p.topic
}


func (p *Publisher) Pub(worker Worker.AbstractWorker, queue chan Cache.Block, RwMutex *sync.RWMutex){
	worker.Produce(queue,RwMutex)
}
//Publisher启动器
func (p *Publisher) Start(listener net.Listener, queue chan Cache.Block, sign chan bool, RwMutex *sync.RWMutex){
	var wpool = new(Worker.Workerpool)
	defer listener.Close()
	wpool.Init(1000000)
	//通道关闭则退出Puber，收到关闭信号量则退出Puber。若出现未知错误，通道均阻塞，则退出Puber
	select {
	case v, ok :=<- queue:
		if !ok{
			return
		}else{
			queue <- v
			for{
				conn, err := listener.Accept()
				Utils.ChkError(err)
				if err != nil{
					return
				}
				go Utils.Handler(conn, wpool)
				for x:= range wpool.GetChannel(){
					worker := x
					go p.Pub(worker, queue, RwMutex)
				}
			}
		}
	case s :=<- sign:
		if s{
			return
		}
	default:
		os.Exit(3)
	}

}



