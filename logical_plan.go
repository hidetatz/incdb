package main

import "sort"

func plan(stmt *QueryStmt) *QueryPlan {
	whr, odr, lim, ofs := stmt.Select.Where, stmt.Select.Order, stmt.Select.Limit, stmt.Select.Offset
	ops := []Operation{
		func(tpls []*Tuple) ([]*Tuple, error) {
			dat, err := readAll()
			if err != nil {
				return nil, err
			}

			r := []*Tuple{}
			for _, d := range dat {
				r = append(r, &Tuple{tp: d})
			}

			return r, nil
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

type Tuple struct {
	tp map[string]string
}

// Operation represents a relational algebra operator.
type Operation func(tpls []*Tuple) ([]*Tuple, error)

func OpWhere(key string) func(tpls []*Tuple) ([]*Tuple, error) {
	return func(tpls []*Tuple) ([]*Tuple, error) {
		for _, tpl := range tpls {
			if _, ok := tpl.tp[key]; ok {
				return []*Tuple{tpl}, nil
			}
		}
		return []*Tuple{}, nil
	}
}

func OpOrder(dir string) func(tpls []*Tuple) ([]*Tuple, error) {
	return func(tpls []*Tuple) ([]*Tuple, error) {
		if dir == "asc" {
			sort.Slice(tpls, func(i, j int) bool {
				var ik string
				for k := range tpls[i].tp {
					ik = k
				}

				var jk string
				for k := range tpls[j].tp {
					jk = k
				}
				return ik < jk
			})

			return tpls, nil
		}
		sort.Slice(tpls, func(i, j int) bool {
			var ik string
			for k := range tpls[i].tp {
				ik = k
			}

			var jk string
			for k := range tpls[j].tp {
				jk = k
			}
			return jk < ik
		})

		return tpls, nil
	}
}

func OpLimitOffset(limit, offset int) func(tpls []*Tuple) ([]*Tuple, error) {
	return func(tpls []*Tuple) ([]*Tuple, error) {
		if limit < 0 {
			limit = len(tpls) // in case limit is negative, it means limit is not specified.
		}

		if limit == 0 {
			return tpls[:0], nil
		}

		if offset > len(tpls) {
			return tpls[:0], nil
		}

		end := offset + limit
		if end > len(tpls) {
			end = len(tpls)
		}

		return tpls[offset:end], nil
	}
}
