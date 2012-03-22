package main

import (
	"strconv"
)

type token int

const (
	INVALID token = iota
	EOF
	NL
	IDENT // foo
	INT // 123
	STRING // "abc"
	TSTART // :xxx{
	TCONT // }xxx
	TEST // ~
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %
	EQ // ==
	LT // <
	GT // >
	ASSIGN // =
	NOT // !
	NE // !=
	LE // <=
	GE // >=
	DOT // .
	OR // ||
	AND // &&
	ntoken
)

var tokname = map[token]string {
	INVALID: "INVALID",
	EOF: "EOF",
	NL: "NL",
	IDENT: "IDENT",
	INT: "INT",
	STRING: "STRING",
	TSTART: "TSTART",
	TCONT: "TCONT",
	TEST: "TEST",
	ADD: "ADD",
	SUB: "SUB",
	MUL: "MUL",
	QUO: "QUO",
	REM: "REM",
	EQ: "EQ",
	LT: "LT",
	GT: "GT",
	ASSIGN: "ASSIGN",
	NOT: "NOT",
	NE: "NE",
	LE: "LE",
	GE: "GE",
	DOT: "DOT",
}

func (t token) String() string {
	name, ok := tokname[t]
	if ok {
		return name
	}
	return "token-" + strconv.Itoa(int(t))
}