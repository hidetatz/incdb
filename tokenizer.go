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

	TkCreate
	TkTable

	TkStr // arbitrary string

	TkLimit
	TkOffset

	TkOrder
	TkBy
	TkAsc
	TkDesc

	TkString // "string" (data type)

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
		case "r":
			cur.Next = &Token{Type: TkRead}
		case "w":
			cur.Next = &Token{Type: TkWrite}
		case "limit":
			cur.Next = &Token{Type: TkLimit}
		case "offset":
			cur.Next = &Token{Type: TkOffset}
		case "order":
			cur.Next = &Token{Type: TkOrder}
		case "by":
			cur.Next = &Token{Type: TkBy}
		case "asc":
			cur.Next = &Token{Type: TkAsc}
		case "desc":
			cur.Next = &Token{Type: TkDesc}
		case "create":
			cur.Next = &Token{Type: TkCreate}
		case "table":
			cur.Next = &Token{Type: TkTable}
		default:
			cur.Next = &Token{Type: TkStr, Val: s}
		}

		cur = cur.Next
	}

	cur.Next = &Token{Type: TkEOF}
	return tk.Next
}
