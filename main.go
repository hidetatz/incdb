package main

import (
	"fmt"
	"os"
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

	if stmt.Insert != nil {
		if err := save(stmt.Insert.Key, stmt.Insert.Val); err != nil {
			fmt.Fprintf(os.Stderr, "save data %s: %s\n", stmt.Insert.Key, err)
			os.Exit(1)
		}
		return
	}

	pln := plan(stmt)

	var result []*Tuple
	for _, ops := range pln.Ops {
		result, _ = ops(result)
	}

	ret := []map[string]string{}
	for _, t := range result {
		ret = append(ret, t.tp)
	}

	fmt.Println(ret)
}
