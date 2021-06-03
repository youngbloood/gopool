package gopool

import "time"

type Pooler interface {
	Send(interface{}, ...time.Duration) error
	Expand(int)
	Cap() int
	Wait()
	Stop()
	Reset()
}
