package main

import (
	"fmt"
	"strconv"
)

// tk is a global token which is "currently" focused on.
// For better parser readability this is globally declared.
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

// "r" str? ("limit" num ("offset" num)?)?
func parseRead() *QueryStmt {
	q := &QueryStmt{Select: &Select{}}

	if consume(TkEOF) {
		return q
	}

	s, ok := consumeStr()
	if ok {
		q.Select.Where = &Where{Equal: &Binary{Value: s}}
	}

	// limit and offset order does not matter (postgres compatible)
	if consume(TkLimit) {
		lim := expectNum()
		if lim < 0 {
			panic("limit must not be negative")
		}
		q.Select.Limit = &Limit{Count: lim}

		if consume(TkOffset) {
			ofs := expectNum()
			if ofs < 0 {
				panic("offset must not be negative")
			}
			q.Select.Offset = &Offset{Count: ofs}
		}
	} else if consume(TkOffset) {
		ofs := expectNum()
		if ofs < 0 {
			panic("offset must not be negative")
		}
		q.Select.Offset = &Offset{Count: ofs}

		if consume(TkLimit) {
			lim := expectNum()
			if lim < 0 {
				panic("limit must not be negative")
			}
			q.Select.Limit = &Limit{Count: lim}
		}
	}

	return q
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
