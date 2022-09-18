package main

// This should have "Left" field when table schema is supported
type Binary struct {
	Value string
}

type Where struct {
	Equal *Binary
}

type Limit struct {
	Count int
}

type Offset struct {
	Count int
}

type Select struct {
	Where  *Where
	Limit  *Limit
	Offset *Offset
}

type Insert struct {
	Key string
	Val string
}

type QueryStmt struct {
	Select *Select
	Insert *Insert
}
