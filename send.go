package main

import (
	"io"
	"log"
)

import "fmt"

type send struct {
	x expr
}

func (s send) eval(e *env) val {
	v := s.x.eval(e).String()
	_, err := io.WriteString(e.conn, v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(">%q\n", v)
	return str(v)
}
