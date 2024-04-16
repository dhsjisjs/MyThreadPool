package Worker

import (
	"sync"
	"sync/atomic"
)

var PayloadPool *sync.Pool
var Created int64

func InitPayloadPool() {
	PayloadPool = &sync.Pool{
		New: func() any {
			atomic.AddInt64(&Created, 1)
			p := new(Payload)
			p.Wait = make(chan any, 1)
			return p
		},
	}
}

type HandleFunc func() any


type Worker struct {
	// Todo: Change it to "Linked List", there is no need to use chanel, we only need it to be concurrent safe
	Payload chan *Payload
}

type Payload struct {
	Do   HandleFunc
	Wait chan any
}

// this will be called everytime a worker is generated
func (w *Worker) run() {
	go func() {
		for {
			Payload := <-w.Payload
			Payload.Wait <- Payload.Do()
		}
	}()
}

type Pool struct {
	pool []*Worker
	// max capacity of workers
	MaxCapacity int
	// which worker will be picker when calling Get(), before it reach maxium capacity, it is always 0
	CurrentWorkerOffset int
	// how many workers are in the pool
	CurrentWorker int
	m             sync.Mutex
}

// returns a new worker that has limit of 2 << 10 *1024, this is quite large
func NewWorker() *Worker {
	w := &Worker{
		//Max Capacity of Payload channel
		Payload: make(chan *Payload, 2<<10*1024),
	}
	w.run()
	return w
}

// returns the worker that handles the payload,  ignore it.
func (p *Pool) Add(Payload *Payload) *Worker {
	if Payload == nil || Payload.Do == nil {
		return nil
	}
	p.m.Lock()
	w := p.Get()
	w.Payload <- Payload
	p.m.Unlock()
	return w
}

// create a new Payload sync pool and a worker pool
func NewPool(max int) *Pool {
	InitPayloadPool()
	return &Pool{
		pool:        make([]*Worker, 0, max),
		MaxCapacity: max,
	}
}

// get a worker from the worker pool
func (p *Pool) Get() *Worker {
	var w *Worker
	if len(p.pool) != p.MaxCapacity {
		w = NewWorker()
		p.pool = append(p.pool, w)
		p.CurrentWorker++

	} else {
		w = p.pool[p.CurrentWorkerOffset]
		if p.CurrentWorkerOffset != p.MaxCapacity-1 {
			p.CurrentWorkerOffset++
		} else {
			p.CurrentWorkerOffset = 0
		}
	}
	return w
}
