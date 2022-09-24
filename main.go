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

	query := strings.Join(os.Args[1:], " ")
	if err := run(query); err != nil {
		fmt.Fprintf(os.Stderr, "run query: %s\n", err)
		os.Exit(1)
	}
}

func run(query string) error {
	stmt, err := parse(query)
	if err != nil {
		return fmt.Errorf("gramatically invalid: %w", err)
	}

	switch {
	case stmt.Create != nil:
		if err := execCreate(stmt.Create); err != nil {
			return fmt.Errorf("execute create statement: %w", err)
		}

		fmt.Fprintf(os.Stdout, "table %s created", stmt.Create.Table)
	case stmt.Insert != nil:
		if err := execInsert(stmt.Insert); err != nil {
			return fmt.Errorf("execute insert statement: %w", err)
		}

		fmt.Fprintf(os.Stdout, "inserted")
	case stmt.Select != nil:
		results, err := execSelect(stmt.Select)
		if err != nil {
			return fmt.Errorf("execute select statement: %w", err)
		}

		for _, r := range results {
			fmt.Fprintf(os.Stdout, "%v", r.Vals)
		}

		fmt.Fprintf(os.Stdout, "\n")
	}

	return nil
}

func execCreate(c *Create) error {
	if err := addTable(c.Table, c.Cols, c.Types); err != nil {
		return fmt.Errorf("add table %s in catalog: %w", c.Table, err)
	}

	if err := createTable(c.Table); err != nil {
		return fmt.Errorf("create table %s: %w", c.Table, err)
	}

	return nil
}

func execInsert(i *Insert) error {
	if err := save(i.Table, i.Cols, i.Vals); err != nil {
		return fmt.Errorf("save data into %s: %w", i.Table, err)
	}

	return nil
}

func execSelect(s *Select) ([]*Record, error) {
	pln := planSelect(s)

	var result []*Record
	for _, ops := range pln.Ops {
		r, err := ops(result)
		if err != nil {
			return nil, fmt.Errorf("compute select result: %w", err)
		}

		result = r
	}

	return result, nil
}
