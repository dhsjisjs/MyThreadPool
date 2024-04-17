package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type FIFONode struct {
	next     *FIFONode
	previous *FIFONode
	data     any
}

func (f *FIFONode) Next() *FIFONode {
	return f.next
}
func (f *FIFONode) Previous() *FIFONode {
	return f.previous
}

func (f *FIFONode) SetNext(next *FIFONode) {
	f.next = next

}
func (f *FIFONode) SetPrevious(pre *FIFONode) {
	f.previous = pre
}

type FIFO struct {
	root   *FIFONode
	length int64
	lock   sync.Mutex
}

// Push an element
func (f *FIFO) Push(node *FIFONode) {
	f.lock.Lock()
	atomic.AddInt64(&f.length, 1)

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
	f.lock.Unlock()
	m.Done()
}

// Pop an element
func (f *FIFO) Pop() *FIFONode {
	f.lock.Lock()
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
	atomic.AddInt64(&f.length, -1)
	f.lock.Unlock()
	return tmp
}

// Return the last push element
func (f *FIFO) Root() *FIFONode {
	return f.root
}

// Return the element will be poped
func (f *FIFO) Previous() *FIFONode {
	return f.root.Previous()
}

var m sync.WaitGroup

func main() {
	m.Add(10)
	f := FIFO{}
	for i := 0; i < 10; i++ {
		j := i
		go func() {
			tmp := FIFONode{
				data: j,
			}
			f.Push(&tmp)

		}()
	}

	m.Wait()
	fmt.Println(f.length)
	m.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(f.Pop())
			m.Done()
		}()
	}
	m.Wait()
	fmt.Println("all done")

}
