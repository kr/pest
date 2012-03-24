package main

import (
	"fmt"
	"strconv"
)

type parserError struct {
	off int
	msg string
}

type parser struct {
	off int
	tok token
	lit string
	err []parserError
	s   scanner
}

type stmt struct {
	tok token
	e expr
}

func (p *parser) init(src []byte) {
	p.s.init(src)
	p.next()
}

func (p *parser) parse() expr {
	var s seq
	for p.tok != EOF {
		p.skipNL()
		s = append(s, p.stmtParse())
	}
	return s
}

func (p *parser) skipNL() {
	for p.tok == NL {
		p.next()
	}
}

func (p *parser) stmtParse() expr {
	defer p.consume(NL)
	switch p.tok {
	case GT:
		p.next()
		x := p.assignParse()
		return send{x}
	case LT:
		p.next()
		off := p.off
		x := p.assignParse()
		return recv{x, string(p.s.src[off:p.off])}
	case TEST:
		p.next()
		off := p.off
		x := p.assignParse()
		return test{x, string(p.s.src[off:p.off])}
	}
	return p.assignParse()
}

func (p *parser) assignParse() expr {
	off := p.off
	e := p.orParse()
	if p.tok == ASSIGN {
		if id, ok := e.(ident); ok {
			p.next()
			r := p.assignParse()
			e = assignment{id, r}
		} else {
			p.error(off, "not an identifier in assignment")
		}
	}
	return e
}

func (p *parser) orParse() expr {
	e := p.andParse()
	for p.tok == OR {
		p.next()
		r := p.andParse()
		e = or{e, r}
	}
	return e
}

func (p *parser) andParse() expr {
	e := p.compareParse()
	for p.tok == AND {
		p.next()
		r := p.compareParse()
		e = and{e, r}
	}
	return e
}

func (p *parser) compareParse() expr {
	e := p.arithParse()
	switch p.tok {
	case EQ:
		p.next()
		r := p.arithParse()
		e = eq{e, r}
	case NE:
		p.next()
		r := p.arithParse()
		e = not{eq{e, r}}
	case GT:
		p.next()
		r := p.arithParse()
		e = gt{e, r}
	case LT:
		p.next()
		r := p.arithParse()
		e = lt{e, r}
	case GE:
		p.next()
		r := p.arithParse()
		e = not{lt{e, r}}
	case LE:
		p.next()
		r := p.arithParse()
		e = not{gt{e, r}}
	}
	return e
}

func (p *parser) arithParse() expr {
	e := p.termParse()
	for {
		switch p.tok {
		case ADD:
			p.next()
			r := p.termParse()
			e = add{e, r}
		case SUB:
			p.next()
			r := p.termParse()
			e = sub{e, r}
		default:
			return e
		}
	}
	panic("unreached")
}

func (p *parser) termParse() expr {
	e := p.unaryParse()
	for {
		switch p.tok {
		case MUL:
			p.next()
			r := p.unaryParse()
			e = mul{e, r}
		case QUO:
			p.next()
			r := p.unaryParse()
			e = quo{e, r}
		case REM:
			p.next()
			r := p.unaryParse()
			e = rem{e, r}
		default:
			return e
		}
	}
	panic("unreached")
}

func (p *parser) unaryParse() expr {
	switch p.tok {
	case NOT:
		p.next()
		return not{p.atomParse()}
	}
	return p.atomParse()
}

func (p *parser) atomParse() expr {
	off, tok, lit := p.off, p.tok, p.lit
	p.next() // always make progress
	switch tok {
	case IDENT:
		return ident(lit)
	case INT:
		return intconv(lit)
	case DOT:
		return dot{}
	case STRING:
		s, err := strconv.Unquote(lit)
		if err != nil {
			panic(err) // can't happen
		}
		return str(s)
	case TSTART:
		val, cont := tconv(lit)
		var e expr = val
		for cont {
			mid := p.assignParse()
			if p.tok != TCONT {
				p.error(p.off, "unexpected ", p.tok)
				return nil
			}
			val, cont = tconv(p.lit)
			p.next()
			e = concat{concat{e, mid}, val}
		}
		return concat{e, str("\r\n")}
	}
	p.error(off, "unexpected ", tok)
	return nil
}

func tconv(lit string) (val str, cont bool) {
	return str(lit[1:len(lit)-1]), lit[len(lit) - 1] == '{'
}

func intconv(lit string) num {
	n, _ := strconv.Atoi(lit)
	return num(n)
}

func (p *parser) error(off int, a ...interface{}) {
	msg := fmt.Sprint(a...)
	p.err = append(p.err, parserError{off, msg})
}

func (p *parser) next() {
	p.off, p.tok, p.lit = p.s.scan()
}

func (p *parser) consume(t token) {
	if p.tok != t {
		p.error(p.off, "unexpected ", p.tok)
	}
	p.next()
}
