package main

import (
	"strconv"
	"strings"
)

type Token struct {
	Type TkType
	Next *Token
	Val  string
	IVal int // active only if Type is TkInt
}

type TkType string

const (
	// Select
	TkSelect = TkType("select")
	TkFrom   = TkType("from")
	TkWhere  = TkType("where")

	TkLimit  = TkType("limit")
	TkOffset = TkType("offset")

	TkOrder = TkType("order")
	TkBy    = TkType("by")
	TkAsc   = TkType("asc")
	TkDesc  = TkType("desc")

	// Insert
	TkInsert = TkType("insert")
	TkInto   = TkType("into")
	TkValues = TkType("values")

	// Create
	TkCreate = TkType("create")
	TkTable  = TkType("table")

	// Data types
	TkString = TkType("string")

	// arbitrary string but not surrounded by quote (e.g. table, column)
	TkSymbol = TkType("symbol")

	TkStr = TkType("string value")
	TkInt = TkType("integer value")

	// Symbols
	TkLParen      = TkType("(")
	TkRParen      = TkType(")")
	TkComma       = TkType(",")
	TkEqual       = TkType("=")
	TkNotEqual    = TkType("!=")
	TkExclamation = TkType("!")
	TkStar        = TkType("*")

	TkEOF = TkType("EOF")
)

func isNumber(b byte) bool {
	return '0' <= b && b <= '9'
}

func isAlphabet(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}

func tokenize(query string) *Token {
	tk := &Token{Next: nil}
	cur := tk

	i := 0
	for i < len(query) {
		switch query[i] {
		case ' ':
			i++
			continue

		case '"':
			i++
			s := ""
			terminated := false
			for i < len(query) {
				if query[i] == '"' { // TODO: must handle " in ""
					terminated = true
					break
				}
				s += string(query[i])
				i++
			}
			if !terminated {
				panic("double quote not terminated")
			}
			i++
			cur.Next = &Token{Type: TkStr, Val: s}

		case '\'':
			i++
			s := ""
			terminated := false
			for i < len(query) {
				if query[i] == '\'' { // TODO: must handle ' in ''
					terminated = true
					break
				}
				s += string(query[i])
				i++
			}
			if !terminated {
				panic("single quote not terminated")
			}
			i++
			cur.Next = &Token{Type: TkStr, Val: s}

		case '(':
			i++
			cur.Next = &Token{Type: TkLParen}

		case ')':
			i++
			cur.Next = &Token{Type: TkRParen}

		case ',':
			i++
			cur.Next = &Token{Type: TkComma}

		case '=':
			i++
			cur.Next = &Token{Type: TkEqual}

		case '!':
			i++
			if query[i] == '=' {
				i++
				cur.Next = &Token{Type: TkNotEqual}
			} else {
				cur.Next = &Token{Type: TkExclamation}
			}
		case '*':
			i++
			cur.Next = &Token{Type: TkStar}

		default:
			s := ""
			for i < len(query) {
				// Some RDB allows using symbol characters in column/table name,
				// but incdb does not to reduce implementation complexity.
				if !isAlphabet(query[i]) && !isNumber(query[i]) {
					break
				}
				s += string(query[i])
				i++
			}

			if i, err := strconv.ParseInt(s, 10, 64); err == nil {
				cur.Next = &Token{Type: TkInt, IVal: int(i)}
				break
			}

			switch strings.ToLower(s) {
			case "select":
				cur.Next = &Token{Type: TkSelect}
			case "from":
				cur.Next = &Token{Type: TkFrom}
			case "where":
				cur.Next = &Token{Type: TkWhere}
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

			case "insert":
				cur.Next = &Token{Type: TkInsert}
			case "into":
				cur.Next = &Token{Type: TkInto}
			case "values":
				cur.Next = &Token{Type: TkValues}

			case "create":
				cur.Next = &Token{Type: TkCreate}
			case "table":
				cur.Next = &Token{Type: TkTable}

			case "string":
				cur.Next = &Token{Type: TkString}

			default:
				cur.Next = &Token{Type: TkSymbol, Val: s}
			}
		}

		cur = cur.Next
	}

	cur.Next = &Token{Type: TkEOF}
	return tk.Next
}
