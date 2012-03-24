package main

import (
	"log"
	"strings"
	"time"
)

type expr interface {
	eval(e *env) val
}

type seq []expr

func (s seq) eval(e *env) val {
	for _, ex := range s {
		ex.eval(e)
	}
	return nil
}


type or struct {
	l, r expr
}

func (o or) eval(e *env) val {
	return nil
}


type and struct {
	l, r expr
}

func (a and) eval(e *env) val {
	return nil
}


type not struct {
	x expr
}

func (n not) eval(e *env) val {
	return boolean(!n.x.eval(e).Bool())
}

type test struct {
	x expr
	src string
}

func (t test) eval(e *env) val {
	if !t.x.eval(e).Bool() {
		log.Fatal("test failed: ", t.src)
	}
	return nil
}

type eq struct {
	l, r expr
}

func (q eq) eval(e *env) val {
	return nil
}


type gt struct {
	l, r expr
}

func (g gt) eval(e *env) val {
	return nil
}


type lt struct {
	l, r expr
}

func (l lt) eval(e *env) val {
	return boolean(l.l.eval(e).Int() < l.r.eval(e).Int())
}


type add struct {
	a, b expr
}

func (a add) eval(e *env) val {
	return nil
}

type sub struct {
	minuend, subtrahend expr
}

func (s sub) eval(e *env) val {
	return num(s.minuend.eval(e).Int() - s.subtrahend.eval(e).Int())
}

type mul struct {
	a, b expr
}

func (m mul) eval(e *env) val {
	return nil
}

type quo struct {
	numerator, denominator expr
}

func (q quo) eval(e *env) val {
	return nil
}

type rem struct {
	numerator, denominator expr
}

func (q rem) eval(e *env) val {
	return nil
}

type ident string

func (i ident) eval(e *env) val {
	switch i {
	case "now":
		return num(int64(time.Now().Sub(time.Time{})))
	}

	v := e.tab[i]
	if v == nil {
		log.Fatal("unbound reference ", i)
	}
	return v
}


type concat struct {
	a, b expr
}

func (c concat) eval(e *env) val {
	a := c.a.eval(e.withDot(nil))
	b := c.b.eval(e.withDot(nil))

	if a != nil && b != nil {
		return str(a.String() + b.String())
	}

	if e.dot == nil {
		return nil
	}

	if a != nil {
		as := a.String()
		if !strings.HasPrefix(*e.dot, as) {
			return nil
		}
		dot := (*e.dot)[len(as):]
		b = c.b.eval(e.withDot(&dot))
		if b == nil {
			return nil
		}
		return str(as + b.String())
	}

	if b != nil {
		bs := b.String()
		if !strings.HasSuffix(*e.dot, bs) {
			return nil
		}
		dot := (*e.dot)[:len(*e.dot)-len(bs)]
		a = c.a.eval(e.withDot(&dot))
		if a == nil {
			return nil
		}
		return str(a.String() + bs)
	}

	// use brute force
	var aok, bok val
	var n int
	for i := range *e.dot {
		adot := (*e.dot)[:i]
		bdot := (*e.dot)[i:]
		a = c.a.eval(e.withDot(&adot))
		b = c.b.eval(e.withDot(&bdot))
		if a != nil && b != nil {
			aok, bok = a, b
			n++
			if n > 1 {
				return nil // ambiguous
			}
		}
	}
	if aok == nil || bok == nil {
		return nil // no match
	}
	return str(aok.String() + bok.String())
}

type assignment struct {
	l ident
	r expr
}

func (a assignment) eval(e *env) val {
	v := a.r.eval(e)
	e.tab[a.l] = v
	return v
}

type dot struct{}

func (d dot) eval(e *env) val {
	if e.dot == nil {
		return nil
	}
	return str(*e.dot)
}
