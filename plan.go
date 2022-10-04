package main

import (
	"sort"
)

func planSelect(slct *Select) *QueryPlan {
	cols, whr, odr, lim, ofs := slct.Columns, slct.Where, slct.Order, slct.Limit, slct.Offset
	ops := []Operation{
		func(rs []*Record) ([]*Record, error) {
			rs, err := readData(slct.Table)
			if err != nil {
				return nil, err
			}

			return rs, nil
		},
	}

	if whr != nil {
		if whr.Equal != nil {
			ops = append(ops, OpWhereEq(whr.Equal.Column, whr.Equal.Value))
		} else {
			ops = append(ops, OpWhereNotEq(whr.NotEqual.Column, whr.NotEqual.Value))
		}
	}

	if odr != nil {
		ops = append(ops, OpOrder(odr.Column, odr.Dir))
	}

	if lim != nil && ofs != nil {
		ops = append(ops, OpLimitOffset(lim.Count, ofs.Count))
	} else if lim != nil {
		ops = append(ops, OpLimitOffset(lim.Count, 0))
	} else if ofs != nil {
		ops = append(ops, OpLimitOffset(-1, ofs.Count))
	}

	ops = append(ops, OpProjection(cols))
	return &QueryPlan{Ops: ops}
}

type QueryPlan struct {
	Ops []Operation
}

// Operation represents a relational algebra operator.
type Operation func(rs []*Record) ([]*Record, error)

func OpWhereEq(col, key string) func(rs []*Record) ([]*Record, error) {
	return func(rs []*Record) ([]*Record, error) {
		i := 0
		for _, r := range rs {
			if r.Find(col, key) {
				rs[i] = r
				i++
			}
		}
		rs = rs[:i]
		return rs, nil
	}
}

func OpWhereNotEq(col, key string) func(rs []*Record) ([]*Record, error) {
	return func(rs []*Record) ([]*Record, error) {
		i := 0
		for _, r := range rs {
			if !r.Find(col, key) {
				rs[i] = r
				i++
			}
		}
		rs = rs[:i]
		return rs, nil
	}
}

func OpOrder(col, dir string) func(rs []*Record) ([]*Record, error) {
	return func(rs []*Record) ([]*Record, error) {
		if len(rs) == 0 {
			return rs, nil
		}

		sort.Slice(rs, func(i, j int) bool {
			if dir == "asc" {
				return rs[i].Value(col) < rs[j].Value(col)
			}
			return rs[j].Value(col) < rs[i].Value(col)
		})
		return rs, nil
	}

}

func OpLimitOffset(limit, offset int) func(rs []*Record) ([]*Record, error) {
	return func(rs []*Record) ([]*Record, error) {
		if len(rs) == 0 {
			return rs, nil
		}

		if limit < 0 {
			limit = len(rs) // in case limit is negative, it means limit is not specified.
		}

		if limit == 0 {
			return rs[:0], nil
		}

		if offset > len(rs) {
			return rs[:0], nil
		}

		end := offset + limit
		if end > len(rs) {
			end = len(rs)
		}

		return rs[offset:end], nil
	}
}

func OpProjection(cols []string) func(rs []*Record) ([]*Record, error) {
	return func(rs []*Record) ([]*Record, error) {
		if cols[0] == "*" {
			return rs, nil
		}

		// First, find which column (index) must be picked up in result
		r := rs[0]
		indices := []int{}
		for i, col := range r.Cols {
			for _, targetCol := range cols {
				if col == targetCol {
					indices = append(indices, i)
					break
				}
			}
		}

		// Then, filter the value by the picked up index
		for _, r := range rs {
			n := 0
			for i := range r.Cols {
				if Contains(indices, i) {
					r.Cols[n] = r.Cols[i]
					r.Types[n] = r.Types[i]
					r.Vals[n] = r.Vals[i]
					n++
				}
			}
			r.Cols = r.Cols[:n]
			r.Types = r.Types[:n]
			r.Vals = r.Vals[:n]
		}

		return rs, nil
	}
}
