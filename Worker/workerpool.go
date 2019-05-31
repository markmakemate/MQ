package Worker


//等候队列
type WaitRoom struct {
	WaitChannel chan AbstractWorker
}
//线程池
type Workerpool struct{
	max_num int
	WorkerChannel chan AbstractWorker
	waitRoom WaitRoom
}
//返回线程池的线程数
func (wp * Workerpool) Size() int{
	return wp.max_num
}
//添加新的worker
func (wp *Workerpool) Push(w AbstractWorker){
	w.SetId(wp.max_num + 1)
	//线程池满时，将Worker转至waitRoom等待,waitRoom满时阻塞, WorkerChannel一旦释放空闲位，waitRoom便推送Worker至WorkerChannel
	select{
	case wp.WorkerChannel <- w:
		wp.max_num = len(wp.WorkerChannel)
	case v :=<- wp.waitRoom.WaitChannel:
		wp.WorkerChannel <- v
	default:
		wp.waitRoom.WaitChannel <- w
	}
}

//线程池初始化
func (wp *Workerpool) Init(MAXSIZE int64){
	wp.WorkerChannel = make(chan AbstractWorker, MAXSIZE)
	wp.waitRoom.WaitChannel = make(chan AbstractWorker, MAXSIZE)
	wp.max_num = 0
}

func (wp *Workerpool) GetChannel() chan AbstractWorker{
	return wp.WorkerChannel
}

