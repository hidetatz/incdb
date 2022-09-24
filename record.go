package main

type Record struct {
	Cols  []string
	Types []string
	Vals  []string
}

func (r *Record) Key() string {
	return r.Vals[0] // assuming first column is key
}

func (r *Record) Find(key string) bool {
	return r.Key() == key // assuming first column is key
}
