package main

import (
	"encoding/json"
	"fmt"
)

type QueryStmt struct {
	Select *Select
	Insert *Insert
	Create *Create
}

func (qs *QueryStmt) String() string {
	if qs == nil {
		return "<nil>"
	}

	// Use json because pretty printing pointer struct is a hassle.
	b, err := json.Marshal(qs)
	if err != nil {
		return fmt.Sprintf("failed to pretty print QueryStmt: %s", err)
	}

	return string(b)
}

/*
 * Select
 */
type Binary struct {
	Column string
	Value  string
}

type Where struct {
	Equal    *Binary
	NotEqual *Binary
}

type Order struct {
	Column string
	Dir    string // asc/desc
}

type Limit struct {
	Count int
}

type Offset struct {
	Count int
}

type Select struct {
	Columns []string
	Table   string
	Where   *Where
	Order   *Order
	Limit   *Limit
	Offset  *Offset
}

/*
 * Insert
 */
type Insert struct {
	Table string
	Cols  []string
	Vals  []string
}

/*
 * Create
 */
type Create struct {
	Table string
	Cols  []string
	Types []string
}
