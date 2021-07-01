package gopool

import (
	"errors"
	"log"
)

var (
	ErrTimeOut = errors.New("send timeout!")
	ErrNotIdle = errors.New("no idle goroutine!")
	ErrStoped  = errors.New("gopool stoped!")
	errHandle  = errors.New("no function to handle the value v!")

	hook HookError = &_hookError{}
)

type HookError interface {
	HookErr(error)
}

func SetHook(he HookError) {
	hook = he
}

type _hookError struct{}

func (_hookError) HookErr(err error) {
	if err != nil {
		log.Printf("[HookError]:%+v\n", err)
	}
}
