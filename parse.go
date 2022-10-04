package main

import (
	"fmt"
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

	if _, ok := consume(TkSelect); ok {
		return parseSelect(), nil
	}

	if _, ok := consume(TkInsert); ok {
		return parseInsert(), nil
	}

	if _, ok := consume(TkCreate); ok {
		return parseCreate(), nil
	}

	return nil, fmt.Errorf("unknown token type: %v", tk.Type)
}

// select = "select" ("*" | columns) "from" table_name where_clause order_clause limit_clause
func parseSelect() *QueryStmt {
	q := &QueryStmt{Select: &Select{}}

	if _, ok := consume(TkStar); ok {
		q.Select.Columns = []string{"*"}
	} else {
		i := 1
		cols := []string{}
		for {
			if i > 100 {
				panic("number of columns must be less than 100")
			}

			s := mustConsume(TkSymbol)
			cols = append(cols, s)

			if _, ok := consume(TkComma); !ok {
				break
			}
			i++
		}
		q.Select.Columns = cols
	}

	mustConsume(TkFrom)
	q.Select.Table = parseTableNameClause()
	q.Select.Where = parseWhereClause()
	q.Select.Order = parseOrderClause()
	q.Select.Limit, q.Select.Offset = parseLimitOffsetClause()

	return q
}

// table_name = symbol
func parseTableNameClause() string {
	return mustConsume(TkSymbol)
}

// where_clause = "where" column_name ("=" | "!=") str
func parseWhereClause() *Where {
	if _, ok := consume(TkWhere); !ok {
		return nil
	}

	w := &Where{}
	col := mustConsume(TkSymbol)

	eq := true
	if _, ok := consume(TkEqual); ok {
		w.Equal = &Binary{Column: col}
	} else if _, ok := consume(TkNotEqual); ok {
		w.NotEqual = &Binary{Column: col}
		eq = false
	} else {
		panic("= or != must be spceified in where clause")
	}

	if eq {
		w.Equal.Value = mustConsume(TkStr)
	} else {
		w.NotEqual.Value = mustConsume(TkStr)
	}

	return w
}

// order_clause = ("order" "by" ("asc" | "desc"))?
func parseOrderClause() *Order {
	if _, ok := consume(TkOrder); !ok {
		return nil
	}

	if _, ok := consume(TkBy); !ok {
		panic("'by' must follow 'order'")
	}

	if _, ok := consume(TkAsc); ok {
		return &Order{Dir: "asc"}
	} else if _, ok := consume(TkDesc); ok {
		return &Order{Dir: "desc"}
	}

	panic("invalid direction specified after 'order by'")
}

// limit_clause = ("limit" num | "offset" num | "limit" num "offset" num | "offset" num "limit" num)?
func parseLimitOffsetClause() (*Limit, *Offset) {
	// limit and offset order does not matter (postgres compatible)

	// Limit comes first
	if _, ok := consume(TkLimit); ok {
		lim := mustConsumeInt()
		if lim < 0 {
			panic("limit must not be negative")
		}
		l := &Limit{Count: lim}

		if _, ok := consume(TkOffset); ok {
			ofs := mustConsumeInt()
			if ofs < 0 {
				panic("offset must not be negative")
			}
			//o.Count = ofs
			return l, &Offset{Count: ofs}
		}

		return l, nil
	}

	// Offset comes first
	if _, ok := consume(TkOffset); ok {
		ofs := mustConsumeInt()
		if ofs < 0 {
			panic("offset must not be negative")
		}
		o := &Offset{Count: ofs}

		if _, ok := consume(TkLimit); ok {
			lim := mustConsumeInt()
			if lim < 0 {
				panic("limit must not be negative")
			}
			return &Limit{Count: lim}, o
		}

		return nil, o
	}

	return nil, nil
}

// "insert" "into" table_name_clause cols? "values" values
func parseInsert() *QueryStmt {
	q := &QueryStmt{Insert: &Insert{}}
	mustConsume(TkInto)

	q.Insert.Table = parseTableNameClause()

	q.Insert.Cols = parseCols()

	mustConsume(TkValues)

	q.Insert.Vals = parseValues()

	return q
}

// cols = "(" col1 "," col2 "," ... ")"
func parseCols() []string {
	i := 1
	ret := []string{}

	if _, ok := consume(TkLParen); !ok {
		return ret // cols is optional
	}

	for {
		if i > 100 {
			panic("cols must be less than 100")
		}

		s := mustConsume(TkSymbol)
		ret = append(ret, s)

		if _, ok := consume(TkRParen); ok {
			break
		}

		mustConsume(TkComma)
		i++
	}

	return ret
}

// values = "(" '"' val1 '"' "," '"' val2 '"' "," ... ")"
func parseValues() []string {
	i := 1
	ret := []string{}

	mustConsume(TkLParen)
	for {
		if i > 100 {
			panic("cols must be less than 100")
		}

		s := mustConsume(TkStr)
		ret = append(ret, s)

		if _, ok := consume(TkRParen); ok {
			break
		}

		mustConsume(TkComma)
		i++
	}

	return ret
}

// "create" "table" table_name_clause "(" column1 "string" "," column2 "string" "," ... ")"
func parseCreate() *QueryStmt {
	q := &QueryStmt{Create: &Create{}}

	mustConsume(TkTable)

	q.Create.Table = parseTableNameClause()

	mustConsume(TkLParen)

	i := 1
	for {
		if i > 100 {
			panic("a table can contain 100 columns at most")
		}

		s := mustConsume(TkSymbol)
		q.Create.Cols = append(q.Create.Cols, s)

		mustConsume(TkString)
		q.Create.Types = append(q.Create.Types, "string")

		if _, ok := consume(TkRParen); ok {
			break
		}

		mustConsume(TkComma)
		i++
	}

	return q
}

func consume(typ TkType) (string, bool) {
	if tk.Type != typ {
		return "", false
	}

	s := tk.Val
	tk = tk.Next
	return s, true
}

func mustConsume(typ TkType) string {
	s, ok := consume(typ)
	if !ok {
		panic(fmt.Sprintf("%s is expected but got %s", string(typ), string(tk.Type)))
	}

	return s
}

func mustConsumeInt() int {
	i := tk.IVal
	mustConsume(TkInt)
	return i
}
