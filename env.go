package main

import (
	"net"
)

type env struct {
	conn net.Conn
	tab  map[ident]val
	dot  *string
}

func (e *env) withDot(s *string) (n *env) {
	n = new(env)
	*n = *e
	n.dot = s
	return n
}
