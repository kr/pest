package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func fmtoff(w io.Writer, src []byte, off int) {
	var line, col int
	for i, c := range src {
		if i >= off {
			break
		}
		col++
		if c == '\n' {
			line++
			col = 0
		}
	}
	fmt.Fprintf(w, "%d:%d", line+1, col+1)
}

func fmterr(w io.Writer, file string, src []byte, off int, msg string) {
	fmt.Fprint(w, file, ":")
	fmtoff(w, src, off)
	fmt.Fprintln(w, ":", msg)
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("usage: pest addr file...")
		return
	}

	addr := args[0]

	tests := make([]expr, len(args[1:]))
	for i, path := range args[1:] {
		tests[i] = mustRead(path)
	}

	for _, t := range tests {
		run(addr, t)
	}
}

func run(addr string, t expr) {
	var e env
	e.tab = make(map[ident]val)
	e.conn = mustDial(addr)
	defer e.conn.Close()
	t.eval(&e)
}

func mustDial(addr string) net.Conn {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func mustRead(path string) expr {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	src, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	var p parser
	p.init(src)
	v := p.parse()
	for _, e := range p.s.err {
		fmterr(os.Stdout, path, src, e.off, e.msg)
	}
	if p.s.err != nil {
		log.Fatal("lexical errors")
	}
	for _, e := range p.err {
		fmterr(os.Stdout, path, src, e.off, e.msg)
	}
	if p.err != nil {
		log.Fatal("syntax errors")
	}

	return v
}
