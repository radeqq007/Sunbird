package token

import "fmt"

type Position struct {
	Filename string
	Offset   int //offset relative to entire file
	Line     int
	Col      int
}

func (p Position) IsValid() bool {
	return p.Line > 0
}

func (p Position) String() string {
	var msg string
	if p.Filename == "" {
		msg = fmt.Sprint(" <", p.Line, ":", p.Col, "> ")
	} else {
		msg = fmt.Sprint(" <", p.Filename, ":", p.Line, ":", p.Col, "> ")
	}

	return msg
}

func (p Position) LineString() string {
	var msg string
	if p.Filename == "" {
		msg = fmt.Sprint(p.Line)
	} else {
		msg = fmt.Sprint(" <", p.Filename, ":", p.Line, "> ")
	}
	return msg
}
