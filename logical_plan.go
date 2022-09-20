package main

import "sort"

type Plan struct {
	Ops []*Operation
}

type Tuple struct {
	tp map[string]string
}

// Operation represents a relational algebra operator.
type Operation struct {
	Fn func(tpls []*Tuple) []*Tuple
}

func OpOrderBy(dir string) func(tpls []*Tuple) []*Tuple {
	return func(tpls []*Tuple) []*Tuple {
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

			return tpls
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

		return tpls
	}
}
