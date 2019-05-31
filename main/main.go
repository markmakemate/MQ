package main

import "MQ"

func main(){
	var q = new(MQ.MQ)
	q.Start("test", "127.0.0.1:8080", "127.0.0.1:8080", 5000)
}
