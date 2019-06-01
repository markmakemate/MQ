package main

import "MQ"

func main(){
	var q = new(MQ.MQ)
	PuberSignal :=make(chan bool)
	SuberSignal :=make(chan bool)
	q.Start("test", "127.0.0.1:8080", "127.0.0.1:8080", 5000, PuberSignal, SuberSignal)
}
