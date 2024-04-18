package FIFO

import (
	"sync"
)

// Node for FIFO
type FIFONode struct {
	next     *FIFONode
	previous *FIFONode
	data     any
}

// Use NewFIFO, or you have to init the sync.Pool maunally
type FIFO struct {
	//Root node of the duoble linked list
	root *FIFONode
	//length of the duoble linked list
	length int64
	// Mutex of Pop and Push operations
	lock sync.Mutex
	//Sync.pool for FIFONode
	pool *sync.Pool
	//noCopy
	noCopy
}

// noCopy
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// Do not use, it is not concurrent safe
func (f *FIFONode) Next() *FIFONode {
	return f.next
}

// Do not use, it is not concurrent safe
func (f *FIFONode) Previous() *FIFONode {
	return f.previous
}

// Do not use, it is not concurrent safe
func (f *FIFONode) SetNext(next *FIFONode) {
	f.next = next

}

// Do not use, it is not concurrent safe
func (f *FIFONode) SetPrevious(pre *FIFONode) {
	f.previous = pre
}

// Use Pop instead
func (f *FIFONode) GetData() any {
	return f.data
}

// Use Push instead
func (f *FIFONode) SetData(data any) {
	f.data = data
}

// Use Pop instead
func (f *FIFO) pop() *FIFONode {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.length == 0 {
		return nil
	}
	tmp := f.root.Previous()
	if f.length == 1 {
		f.root = nil
	} else {
		tmp.Previous().SetNext(f.root)
		f.root.SetPrevious(tmp.Previous())
	}
	f.length--
	return tmp
}

// Use Push instead
func (f *FIFO) push(node *FIFONode) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.length++

	if f.root != nil {
		//Push the node into the root
		node.SetNext(f.root)
		node.SetPrevious(f.root.previous)
		f.root = node
		node.next.SetPrevious(node)
	} else {
		f.root = node
		f.root.SetNext(f.root)
		f.root.SetPrevious(f.root)
	}
}

// Push an element
func (f *FIFO) Push(data any) {
	var p = f.pool.Get().(*FIFONode)
	p.SetData(data)
	f.push(p)
}

// Pop an element if there is one, else it will return nil
func (f *FIFO) Pop() any {
	var p = f.pop()
	if p != nil {
		tmp := p.GetData()
		p.SetData(nil)
		f.pool.Put(p)
		return tmp
	}
	return nil
}

// Return the head element
func (f *FIFO) Head() *FIFONode {
	return f.root
}

// Return the tail element
func (f *FIFO) Previous() *FIFONode {
	return f.root.Previous()
}

// Return the length of fifo
func (f *FIFO) GetLength() int {
	return int(f.length)
}

// Create an fifo
func NewFIFO() *FIFO {
	var f FIFO
	f.pool = &sync.Pool{
		New: func() any {
			return new(FIFONode)
		},
	}
	return &f
}
