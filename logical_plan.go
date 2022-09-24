package main

import (
	"sort"
)

func plan(stmt *QueryStmt) *QueryPlan {
	whr, odr, lim, ofs := stmt.Select.Where, stmt.Select.Order, stmt.Select.Limit, stmt.Select.Offset
	ops := []Operation{
		func(rs []*Record) ([]*Record, error) {
			rs, err := readData(stmt.Select.Table)
			if err != nil {
				return nil, err
			}

			return rs, nil
		},
	}

	if whr != nil {
		ops = append(ops, OpWhere(whr.Equal.Value))
	}

	if odr != nil {
		ops = append(ops, OpOrder(odr.Dir))
	}

	if lim != nil && ofs != nil {
		ops = append(ops, OpLimitOffset(lim.Count, ofs.Count))
	} else if lim != nil {
		ops = append(ops, OpLimitOffset(lim.Count, 0))
	} else if ofs != nil {
		ops = append(ops, OpLimitOffset(-1, ofs.Count))
	}
	return &QueryPlan{Ops: ops}
}

type QueryPlan struct {
	Ops []Operation
}

// Operation represents a relational algebra operator.
type Operation func(rs []*Record) ([]*Record, error)

func OpWhere(key string) func(rs []*Record) ([]*Record, error) {
	return func(rs []*Record) ([]*Record, error) {
		for _, r := range rs {
			if r.Find(key) {
				return []*Record{r}, nil
			}
		}

		return []*Record{}, nil
	}
}

func OpOrder(dir string) func(rs []*Record) ([]*Record, error) {
	return func(rs []*Record) ([]*Record, error) {
		sort.Slice(rs, func(i, j int) bool {
			if dir == "asc" {
				return rs[i].Key() < rs[j].Key()
			}
			return rs[j].Key() < rs[i].Key()
		})
		return rs, nil
	}

}

func OpLimitOffset(limit, offset int) func(rs []*Record) ([]*Record, error) {
	return func(rs []*Record) ([]*Record, error) {
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
