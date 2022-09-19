package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "at least one argument is required\n")
		os.Exit(1)
	}

	stmt, err := parse(strings.Join(os.Args[1:], " "))
	if err != nil {
		fmt.Fprintf(os.Stderr, "gramatically invalid: %s\n", err)
		os.Exit(1)
	}

	switch {
	case stmt.Select != nil:
		if stmt.Select.Where != nil {
			out, err := readAll()
			if err != nil {
				fmt.Fprintf(os.Stderr, "read data: %s\n", err)
				os.Exit(1)
			}

			for _, m := range out {
				v, ok := m[stmt.Select.Where.Equal.Value]
				if !ok {
					continue
				}

				fmt.Println(v)
				return
			}

			fmt.Println("not found")
			return
		}

		out, err := readAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "read all the data: %s\n", err)
			os.Exit(1)
		}

		if stmt.Select.Order != nil {

			if stmt.Select.Order.Dir == "asc" {
				sort.Slice(out, func(i, j int) bool {
					var ik string
					for k := range out[i] {
						ik = k
					}

					var jk string
					for k := range out[j] {
						jk = k
					}
					return ik < jk
				})
			} else {
				sort.Slice(out, func(i, j int) bool {
					var ik string
					for k := range out[i] {
						ik = k
					}

					var jk string
					for k := range out[j] {
						jk = k
					}
					return jk < ik
				})
			}
		}

		// in case offset is specified, cut the out slice from offset to the last
		if stmt.Select.Offset != nil {
			if stmt.Select.Offset.Count > len(out) {
				stmt.Select.Offset.Count = len(out)
			}

			out = out[stmt.Select.Offset.Count:]
		}

		if stmt.Select.Limit == nil || stmt.Select.Limit.Count > len(out) {
			fmt.Println(out)
			return
		}

		b := []map[string]string{}
		for i := 0; i < stmt.Select.Limit.Count; i++ {
			b = append(b, out[i])
		}
		fmt.Println(b)

	case stmt.Insert != nil:
		if err := save(stmt.Insert.Key, stmt.Insert.Val); err != nil {
			fmt.Fprintf(os.Stderr, "save data %s: %s\n", stmt.Insert.Key, err)
			os.Exit(1)
		}
	}

}
