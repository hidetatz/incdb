package main

import (
	"fmt"
	"strconv"
)

// tk is a global token which is "currently" focused on.
// This is globally declared for better parser readability.
var tk *Token

func parse(query string) (queryStmt *QueryStmt, err error) {
	// For better readability, use panic/recover to get back here from deep-nested parser on error.
	// The argument of panic() will be caught and returned to the caller.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse statement: %v", r)
		}
	}()

	tk = tokenize(query)

	if consume(TkRead) {
		return parseRead(), nil
	}

	if consume(TkWrite) {
		return parseWrite(), nil
	}

	return nil, fmt.Errorf("unknown token type: %v", tk.Type)
}

// read = "r" str? order_clause? limit_clause?
func parseRead() *QueryStmt {
	q := &QueryStmt{Select: &Select{}}

	if consume(TkEOF) {
		return q
	}

	s, ok := consumeStr()
	if ok {
		q.Select.Where = &Where{Equal: &Binary{Value: s}}
	}

	q.Select.Order = parseOrderClause()
	q.Select.Limit, q.Select.Offset = parseLimitOffsetClause()

	return q
}

// order_clause = ("order" "by" ("asc" | "desc"))
func parseOrderClause() *Order {
	if !consume(TkOrder) {
		return nil
	}

	if !consume(TkBy) {
		panic("'by' must follow 'order'")
	}

	if consume(TkAsc) {
		return &Order{Dir: "asc"}
	} else if consume(TkDesc) {
		return &Order{Dir: "desc"}
	}

	panic("invalid direction specified after 'order by'")
}

// limit_clause = ("limit" num | "offset" num | "limit" num "offset" num | "offset" num "limit" num)
func parseLimitOffsetClause() (*Limit, *Offset) {
	// limit and offset order does not matter (postgres compatible)

	// Limit comes first
	if consume(TkLimit) {
		lim := expectNum()
		if lim < 0 {
			panic("limit must not be negative")
		}
		l := &Limit{Count: lim}

		if consume(TkOffset) {
			ofs := expectNum()
			if ofs < 0 {
				panic("offset must not be negative")
			}
			//o.Count = ofs
			return l, &Offset{Count: ofs}
		}

		return l, nil
	}

	// Offset comes first
	if consume(TkOffset) {
		ofs := expectNum()
		if ofs < 0 {
			panic("offset must not be negative")
		}
		o := &Offset{Count: ofs}

		if consume(TkLimit) {
			lim := expectNum()
			if lim < 0 {
				panic("limit must not be negative")
			}
			return &Limit{Count: lim}, o
		}

		return nil, o
	}

	return nil, nil
}

// "w" str str
func parseWrite() *QueryStmt {
	key := mustStr()
	val := mustStr()
	return &QueryStmt{Insert: &Insert{Key: key, Val: val}}
}

func consume(typ TkType) bool {
	if tk.Type == typ {
		tk = tk.Next
		return true
	}

	return false
}

func consumeStr() (string, bool) {
	if tk.Type != TkStr {
		return "", false
	}

	s := tk.Val
	tk = tk.Next
	return s, true
}

func mustStr() string {
	s, ok := consumeStr()
	if !ok {
		panic("not a string")
	}

	return s
}

func expectNum() int {
	s, ok := consumeStr()
	if !ok {
		panic("not a string")
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic("cannot convert to int")
	}

	return int(n)
}
