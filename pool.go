package gopool

import (
	"sync"
	"time"
)

type goPool struct {
	rwMux   sync.RWMutex // 针对queue的Pop()加锁，防止done() pop一个删除的数据时，Send同时发生引发错误
	size    int
	queue   queuer
	run     Run
	isStart bool
	once    sync.Once
	sendTTL *time.Ticker
}

func New(size int, run Run) *goPool {
	return &goPool{
		size:    size,
		run:     run,
		queue:   newQueue(size),
		sendTTL: time.NewTicker(1 * time.Second),
	}
}

// 开启goroutine
func (gp *goPool) StartGo() {
	if gp.isStart {
		return
	}

	gp.once.Do(func() {
		for i := 0; i < gp.size; i++ {
			gochan := gp.queue.Pop()
			go gp.goFunc(gochan)
		}
	})

	gp.isStart = true
}

func (gp *goPool) goFunc(gc *goChan) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				hook.HookErr(err)
			}
		}
	}()

	for v := range gc.dataChan {
		if err := gp.run(v); err != nil {
			if hookErr, ok := v.(HookError); ok {
				hookErr.HookErr(err)
			} else {
				hook.HookErr(err)
			}
		}
	}
}

func (gp *goPool) Send(v interface{}, ttls ...time.Duration) error {
	if !gp.isStart {
		gp.StartGo()
	}

	gp.rwMux.RLock()
	defer gp.rwMux.RUnlock()

	gochan := gp.queue.Pop()
	if gochan == nil {
		return ErrNotIdle
	}
	if len(ttls) > 0 {
		gp.sendTTL.Reset(ttls[0])
	}
	select {
	case gochan.dataChan <- v:
	case <-gp.sendTTL.C:
		return ErrTimeOut
	}
	return nil
}

// 增加容量
func (gp *goPool) Add(x int) {
	if x <= 0 {
		return
	}
	gp.rwMux.Lock()
	defer gp.rwMux.Unlock()

	gp.size += x
	gp.queue.Expan(x, gp.goFunc)
}

// 减少容量
func (gp *goPool) Done(x int) {
	if x <= 0 {
		return
	}
	if x > gp.size {
		x = gp.size
	}
	if x == 0 {
		return
	}

	gp.rwMux.Lock()
	defer gp.rwMux.Unlock()

	gp.size -= x
	gp.queue.Expan(-1*x, gp.goFunc)

	for i := 0; i < x; i++ {
		gochan := gp.queue.Pop()
		close(gochan.dataChan)
	}
}

func (gp *goPool) Wait() {
	for gp.size != 0 {
		<-gp.sendTTL.C
	}
}

func (gp *goPool) Cap() int {
	return gp.size
}
