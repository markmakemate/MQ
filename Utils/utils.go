package Utils

import (
	"MQ/Worker"
	"log"
	"net"
)

//错误日志
func ChkError(err error){
	if err != nil{
		log.Fatal(err)
	}
}
//监听启动器，返回一个Listener对象
func Listen(addr string) net.Listener{
	listener, err := net.Listen("tcp", addr)
	ChkError(err)
	return listener
}
//连接处理器，将连接放入连接池内
func Handler(conn net.Conn, wpool *Worker.Workerpool){
	var w Worker.AbstractWorker
	w = new(Worker.Worker)
	w.SetConn(conn)
	wpool.Push(w)
}
