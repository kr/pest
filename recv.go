package main

import (
	"io"
	"log"
)

import "fmt"

type recv struct {
	x   expr
	src string
}

func (r recv) eval(e *env) val {
	s, err := readLine(e.conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("<%q\n", s)

	v := r.x.eval(e.withDot(&s))
	if v == nil {
		log.Fatalf("ambiguous or no match %#v != %v", s[:len(s)-2], r.src)
	}
	if s != v.String() {
		log.Fatalf("mismatch %#v != %#v", s, v)
	}

	return str(s)
}

func readLine(r io.Reader) (string, error) {
	var buf []byte
	var p byte
	for {
		b := make([]byte, 1)
		_, err := r.Read(b)
		if err != nil {
			return "", err
		}
		buf = append(buf, b[0])
		if p == '\r' && b[0] == '\n' {
			break
		}
		p = b[0]
	}
	return string(buf), nil
}
