package token

import (
	"strconv"
)

type Position struct {
	Filename string
	Offset   int // offset relative to entire file
	Line     int
	Col      int
}

func (p Position) IsValid() bool {
	return p.Line > 0
}

func (p Position) String() string {
	var msg string
	if p.Filename == "" {
		msg = " <" + strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Col) + "> "
	} else {
		msg = " <" + p.Filename + ":" + strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Col) + "> "
	}

	return msg
}

func (p Position) LineString() string {
	var msg string
	if p.Filename == "" {
		msg = strconv.Itoa(p.Line)
	} else {
		msg = " <" + p.Filename + ":" + strconv.Itoa(p.Line) + "> "
	}
	return msg
}
