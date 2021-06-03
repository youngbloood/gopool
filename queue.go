package gopool

import (
	"sync"
)

type node struct {
	val       *goChan
	pre, next *node
}

type goChan struct {
	dataChan chan interface{}
}

type queue struct {
	rwMux            sync.RWMutex
	size, reduceSize int
	head, tail       *node
}

type queuer interface {
	Pop() *goChan
	Expand(int, func(*goChan))
}

func newNode() *node {
	return &node{
		val: &goChan{
			dataChan: make(chan interface{}),
		},
	}
}

func newQueue(size int) *queue {
	q := &queue{
		size: size,
	}
	if size <= 0 {
		return q
	}

	n := newNode()
	q.head = n
	q.tail = n

	for i := 0; i < size-1; i++ {
		n = newNode()
		n.pre = q.tail
		q.tail.next = n
		q.tail = n
	}
	return q
}

func (q *queue) Pop() *goChan {
	q.rwMux.RLock()
	defer q.rwMux.RUnlock()

	if q.head == nil {
		return nil
	}
	if q.head.next == nil && q.reduceSize > 0 {
		val := q.head.val
		q.reset()
		return val
	}
	if q.head.next == nil {
		return q.head.val
	}

	// 1. 移除queue的head
	node := q.head
	node.next.pre = nil
	q.head = node.next
	node.next = nil
	val := node.val

	if q.reduceSize <= 0 {
		// 2. 将head添加到queue的队尾
		q.tail.next = node
		node.pre = q.tail
		q.tail = node
	} else {
		// 3. 丢弃该node
		q.reduceSize--
	}
	return val
}

func (q *queue) Expand(size int, run func(*goChan)) {
	if size == 0 {
		return
	}
	q.rwMux.Lock()
	defer q.rwMux.Unlock()
	q.size += size
	if size > 0 {
		q.add(size, run)
	} else {
		q.done(-1 * size)
	}
}

func (q *queue) add(size int, run func(*goChan)) {
	for i := 0; i < size; i++ {
		gochan := &goChan{
			dataChan: make(chan interface{}),
		}
		n := &node{
			pre: q.tail,
			val: gochan,
		}
		if q.tail != nil {
			q.tail.next = n
			q.tail = n
		} else {
			q.tail = n
			q.head = q.tail
		}

		go run(gochan)
	}
}

func (q *queue) done(size int) {
	// 待pop之后安全移除无数据节点
	q.reduceSize += size
}

func (q *queue) reset() {
	q.head = nil
	q.tail = nil
	q.reduceSize = 0
	q.reduceSize = 0
}
