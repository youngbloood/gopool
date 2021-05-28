package gopool

import "errors"

type HookError interface {
	HookErr(error)
}

var hook HookError

func SetHook(hook HookError) {
	hook = hook
}

var (
	ErrTimeOut = errors.New("send timeout!")
	ErrNotIdle = errors.New("no idle goroutine!")
)
