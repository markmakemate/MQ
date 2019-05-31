package Worker

import (
	"net"
	"sync"
)
import "MQ/Cache"
//抽象worker
type AbstractWorker interface {
	GetConn() net.Conn
	Produce(queue chan Cache.Block, rwMutex *sync.RWMutex)
	Consume(block Cache.Block)
	SetConn(conn net.Conn)
	Close()
	GetId() int
	SetId(id int)
}






