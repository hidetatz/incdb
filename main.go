package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "single query must be passed\n")
		os.Exit(1)
	}

	stmt, err := parse(strings.Join(os.Args[1:], " "))
	if err != nil {
		fmt.Fprintf(os.Stderr, "gramatically invalid: %s\n", err)
		os.Exit(1)
	}

	//Debug(stmt)

	if stmt.Create != nil {
		if err := addTable(stmt.Create.Table, stmt.Create.Cols, stmt.Create.Types); err != nil {
			fmt.Fprintf(os.Stderr, "add table %s in catalog: %s\n", stmt.Create.Table, err)
			os.Exit(1)
		}

		if err := createTable(stmt.Create.Table); err != nil {
			fmt.Fprintf(os.Stderr, "create table %s: %s\n", stmt.Create.Table, err)
			os.Exit(1)
		}
		return
	}

	if stmt.Insert != nil {
		if err := save(stmt.Insert.Table, stmt.Insert.Cols, stmt.Insert.Vals); err != nil {
			fmt.Fprintf(os.Stderr, "save data into %s: %s\n", stmt.Insert.Table, err)
			os.Exit(1)
		}
		return
	}

	pln := plan(stmt)

	var result []*Record
	for _, ops := range pln.Ops {
		r, err := ops(result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "compute result: %s\n", err)
			os.Exit(1)
		}

		result = r
	}

	for _, r := range result {
		fmt.Printf("%v", r.Vals)
	}

	fmt.Println()
}
