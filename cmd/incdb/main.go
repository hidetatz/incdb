package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if err := runQuery(); err != nil {
		fmt.Fprintf(os.Stderr, "incdb: %s\n", err)
		os.Exit(1)
	}
}

func runQuery() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("one argument is required")
	}

	type Req struct {
		Query string
	}

	query := os.Args[1]

	qReq := Req{Query: query}
	b, err := json.Marshal(&qReq)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:2134/query", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("call server: %w", err)
	}

	type Result struct {
		Msg      string
		Hdr      []string
		Vals     [][]string
		ErrorMsg string
	}

	result := Result{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if result.ErrorMsg != "" {
		return fmt.Errorf("run query: %s", result.ErrorMsg)
	}

	if result.Msg != "" {
		fmt.Println(result.Msg)
		return nil
	}

	if os.Getenv("INCDB_TEST") == "1" {
		// in test, output will be structured for testability
		type Output struct {
			Hdr  []string
			Vals [][]string
		}
		o := Output{Hdr: result.Hdr, Vals: result.Vals}

		b, err := json.Marshal(&o)
		if err != nil {
			return fmt.Errorf("marshal result: %w", err)
		}
		fmt.Println(string(b))

		return nil
	}

	tw := NewTableWriter(os.Stdout)
	tw.SetHeader(result.Hdr)
	for _, val := range result.Vals {
		tw.Append(val)
	}

	tw.Render()
	return nil
}
