package gopool

import (
	"sync"
	"time"
)

type goPool struct {
	rwMux           sync.RWMutex // 针对queue的Pop()加锁，防止done() pop一个删除的数据时，Send同时发生引发错误
	size            int
	queue           queuer
	run             Run
	isStart, isStop bool
	once            sync.Once
	sendTTL         *time.Ticker
}

func New(size int, runs ...Run) *goPool {
	var run Run
	if len(runs) > 0 {
		run = runs[0]
	}

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
		var err error
		if run2, ok := v.(Run2); ok {
			err = run2()
		} else if gp.run != nil {
			err = gp.run(v)
		} else {
			hook.HookErr(errHandle)
			continue
		}

		if err != nil {
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

	if gp.isStop {
		return ErrStoped
	}

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

// 增加/减少容量
func (gp *goPool) Expand(x int) {
	if x == 0 {
		return
	}

	gp.rwMux.Lock()
	defer gp.rwMux.Unlock()
	gp.size += x
	if x > 0 {
		gp.queue.Expand(x, gp.goFunc)
		return
	}
	gp.done(-1 * x)
}

// 减少容量
func (gp *goPool) done(x int) {
	if x > gp.size {
		x = gp.size
	}

	gp.size -= x
	gp.queue.Expand(-1*x, gp.goFunc)

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

func (gp *goPool) Stop() {
	gp.rwMux.Lock()
	defer gp.rwMux.Unlock()
	if gp.isStop {
		return
	}
	gp.isStop = true
	gp.done(gp.size)
}

func (gp *goPool) Reset() {
	gp.rwMux.Lock()
	defer gp.rwMux.Unlock()
	gp.Stop()
	gp.size = 0
	gp.isStart = false
	gp.isStop = false
	gp.sendTTL = time.NewTicker(1 * time.Second)
	gp.queue = newQueue(0)
}
