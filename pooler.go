package gopool

import "time"

type Pooler interface {
	Send(interface{}, ...time.Duration) error
	Add(int)
	Done(int)
	Cap() int
	Wait()
}
