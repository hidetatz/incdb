package main

import "fmt"

// r key?
// w key val
func parse(query string) (*QueryStmt, error) {
	tk := tokenize(query)

	if tk.Type == TkRead {
		tk = tk.Next
		if tk.Type == TkEOF {
			return &QueryStmt{Select: &Select{}}, nil
		}

		key := expectStr(tk)
		return &QueryStmt{Select: &Select{Where: &Where{Equal: &Binary{Value: key}}}}, nil
	}

	if tk.Type == TkWrite {
		tk = tk.Next

		key := expectStr(tk)
		tk = tk.Next
		val := expectStr(tk)
		if tk.Type == TkStr {
			return &QueryStmt{Insert: &Insert{Key: key, Val: val}}, nil
		}
	}

	return nil, fmt.Errorf("unknown token type: %v", tk.Type)
}

func expectStr(tk *Token) string {
	if tk.Type != TkStr {
		panic("not a string")
	}

	v := tk.Val
	return v
}
