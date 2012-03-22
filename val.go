package main

import (
	"strconv"
)

type val interface{
	String() string
	Int() int64
	Bool() bool
}

type str string

func (s str) eval(e *env) val {
	return s
}

func (s str) Bool() bool {
	return true
}

func (s str) Int() int64 {
	return 0
}

func (s str) String() string {
	return string(s)
}

type num int64

func (n num) eval(e *env) val {
	return n
}

func (n num) Bool() bool {
	return true
}

func (n num) Int() int64 {
	return int64(n)
}

func (n num) String() string {
	return strconv.FormatInt(int64(n), 10)
}

type boolean bool

func (b boolean) eval(e *env) val {
	return b
}

func (b boolean) Bool() bool {
	return bool(b)
}

func (b boolean) Int() int64 {
	return 0
}

func (b boolean) String() string {
	if b {
		return "true"
	}
	return "false"
}
