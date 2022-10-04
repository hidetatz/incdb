package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func postQuery(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Query string
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("decode request: %s\n", err)
		json.NewEncoder(w).Encode(&Result{ErrorMsg: "internal server error"})
		return
	}

	fmt.Printf("query: %s\n", req.Query)

	result, err := runQuery(req.Query)
	if err != nil {
		fmt.Printf("run query: %s\n", err)
		json.NewEncoder(w).Encode(&Result{ErrorMsg: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(result)
}

type Result struct {
	Msg      string
	Hdr      []string
	Vals     [][]string
	ErrorMsg string
}

func runQuery(query string) (*Result, error) {
	stmt, err := parse(query)
	if err != nil {
		return nil, fmt.Errorf("gramatically invalid: %w", err)
	}

	Debug("statement: ", stmt)

	switch {
	case stmt.Create != nil:
		if err := execCreate(stmt.Create); err != nil {
			return nil, fmt.Errorf("execute create statement: %w", err)
		}

		return &Result{Msg: fmt.Sprintf("table %s created", stmt.Create.Table)}, nil

	case stmt.Insert != nil:
		if err := execInsert(stmt.Insert); err != nil {
			return nil, fmt.Errorf("execute insert statement: %w", err)
		}

		return &Result{Msg: "inserted"}, nil

	case stmt.Select != nil:
		results, err := execSelect(stmt.Select)
		if err != nil {
			return nil, fmt.Errorf("execute select statement: %w", err)
		}

		if len(results) == 0 {
			return &Result{Msg: "no results"}, nil
		}

		res := Result{Hdr: results[0].Cols}

		vals := [][]string{}
		for _, r := range results {
			vals = append(vals, r.Vals)
		}
		res.Vals = vals

		return &res, nil
	}

	panic("never come")
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
