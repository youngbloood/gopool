package gopool

type Run func(interface{}) error

type Run2 = func() error
