package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"test/Worker"
	"time"
)

func sayhello() any {
	time.Sleep(time.Millisecond * 10)
	fmt.Println("hello!")
	return "hello"
}

func main() {
	fmt.Println("http server start")
	err := http.ListenAndServe("localhost:8090", nil)
	if err != nil {
		log.Fatal(err)
	}
	var n = sync.WaitGroup{}
	var p = Worker.NewPool(100)
	n.Add(10000)
	for i := 0; i < 10000; i++ {
		time.Sleep(time.Millisecond)
		go func() {
			var q = Worker.PayloadPool.Get().(*Worker.Payload)
			q.Do = sayhello
			p.Add(q)
			<-q.Wait
			Worker.PayloadPool.Put(q)
			n.Done()
		}()
	}
	n.Wait()
	fmt.Printf("The pool has %v workers\n", p.CurrentWorker)
	fmt.Printf("%v Payloads were created\n", Worker.Created)
	time.Sleep(time.Hour)
}
