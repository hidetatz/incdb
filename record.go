package main

type Record struct {
	Cols  []string
	Types []string
	Vals  []string
}

func (r *Record) ColIndex(col string) int {
	for i := range r.Cols {
		if r.Cols[i] == col {
			return i
		}
	}
	return -1
}

func (r *Record) Value(col string) string {
	index := r.ColIndex(col)
	if index < 0 {
		return ""
	}

	return r.Vals[index]
}

func (r *Record) Find(col, key string) bool {
	index := r.ColIndex(col)
	if index < 0 {
		return false
	}

	return r.Vals[index] == key
}
