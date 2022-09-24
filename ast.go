package main

type QueryStmt struct {
	Select *Select
	Insert *Insert
	Create *Create
}

/*
 * Select
 */

// This should have "Left" field when table schema is supported
type Binary struct {
	Value string
}

type Where struct {
	Equal *Binary
}

type Order struct {
	Dir string // asc/desc
}

type Limit struct {
	Count int
}

type Offset struct {
	Count int
}

type Select struct {
	Table  string
	Where  *Where
	Order  *Order
	Limit  *Limit
	Offset *Offset
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
