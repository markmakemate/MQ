package Worker

import (
	"MQ/Cache"
	"log"
	"net"
	"os"
	"sync"
)
//错误日志
func ChkError(err error){
	if err != nil{
		log.Fatal(err)
	}
}
//调度者
type Worker struct{
	id int
	conn net.Conn
}

//Set, Get方法
func (w *Worker) GetConn() net.Conn{
	return w.conn
}
func (w *Worker) SetConn(c net.Conn){
	w.conn = c
}
func (w * Worker) GetId() int{
	return w.id
}
func (w *Worker) SetId(id int){
	w.id = id
}

//单个worker生产
func (w *Worker) Produce(queue chan Cache.Block, mutex *sync.RWMutex){
	buf := make([]byte, 32)
	var b = new(Cache.Block)
	for{
		mutex.Lock()
		n, err := w.conn.Read(buf)
		if err != nil{
			mutex.Unlock()
			log.Fatal(err)
		}else{
			b.Set(buf[:n])
			b.SetOffset(len(queue))
			//queue关闭时退出
			select {
			case _, ok :=<- queue:
				if !ok{
					mutex.Unlock()
					w.Close()
					return
				}
			case queue <- *b:
				mutex.Unlock()
			default:
				os.Exit(3)
			}
		}
	}

}

//单个worker消费
func (w *Worker) Consume(block Cache.Block){
	_, err := w.GetConn().Write(block.Get())
	ChkError(err)
}

func (w *Worker) Close(){
	_ = w.conn.Close()
}


