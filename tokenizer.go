package main

import "strings"

type Token struct {
	Type TkType
	Next *Token
	Val  string
}

type TkType int

const (
	TkRead TkType = iota + 1
	TkWrite
	TkStr
	TkLimit
	TkEOF
)

func isspace(c byte) bool {
	return c == ' '
}

func tokenize(query string) *Token {
	tk := &Token{Next: nil}
	cur := tk

	i := 0
	for i < len(query) {
		if isspace(query[i]) {
			i++
			continue
		}

		if query[i] == 'r' {
			cur.Next = &Token{Type: TkRead}
			cur = cur.Next
			i++
			continue
		}

		if query[i] == 'w' {
			cur.Next = &Token{Type: TkWrite}
			cur = cur.Next
			i++
			continue
		}

		// arbitrary string
		s := ""
		for i < len(query) {
			if isspace(query[i]) {
				break
			}
			s += string(query[i])
			i++
		}

		switch strings.ToLower(s) {
		case "limit":
			cur.Next = &Token{Type: TkLimit}
			cur = cur.Next
		default:
			cur.Next = &Token{Type: TkStr, Val: s}
			cur = cur.Next
		}
	}

	cur.Next = &Token{Type: TkEOF}
	return tk.Next
}
