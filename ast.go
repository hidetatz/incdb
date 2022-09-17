package main

// This should have "Left" field when table schema is supported
type Binary struct {
	Value string
}

type Where struct {
	Equal *Binary
}

type Select struct {
	Where *Where
}

type Insert struct {
	Key string
	Val string
}

type QueryStmt struct {
	Select *Select
	Insert *Insert
}
