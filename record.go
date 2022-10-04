package main

type Record struct {
	Cols  []string
	Types []string
	Vals  []string
}

func (r *Record) Key() string {
	return r.Vals[0]
}

func (r *Record) Find(col, key string) bool {
	i := 0
	for j, rcol := range r.Cols {
		if rcol == col {
			i = j
		}
	}

	return r.Vals[i] == key
}
