package main

import (
	"reflect"
	"testing"
)

func TestParse_Select(t *testing.T) {
	tests := []struct {
		input    string
		wantStmt *QueryStmt
		wantErr  bool
	}{
		// Currently the parser isn't strict enough to never miss grammer mistakes.
		// In the future it should be improved.
		{
			input:   "r",
			wantErr: true,
		},
		{
			input:    "r test",
			wantStmt: &QueryStmt{Select: &Select{Table: "test"}},
		},
		{
			input:    "r test '1'",
			wantStmt: &QueryStmt{Select: &Select{Table: "test", Where: &Where{Equal: &Binary{Value: "1"}}}},
		},
		{
			input:    `r test "1"`,
			wantStmt: &QueryStmt{Select: &Select{Table: "test", Where: &Where{Equal: &Binary{Value: "1"}}}},
		},
		{
			input:   `r test "1`,
			wantErr: true,
		},
		{
			input:    `r test order by asc`,
			wantStmt: &QueryStmt{Select: &Select{Table: "test", Order: &Order{Dir: "asc"}}},
		},
		{
			input:    `r test order by desc`,
			wantStmt: &QueryStmt{Select: &Select{Table: "test", Order: &Order{Dir: "desc"}}},
		},
		{
			input:   `r test order by invalid`,
			wantErr: true,
		},
		{
			input:   `r test order`,
			wantErr: true,
		},
		{
			input:   `r test order by`,
			wantErr: true,
		},
		{
			input:   `r test order xx asc`,
			wantErr: true,
		},
		{
			input:    `r test limit 100`,
			wantStmt: &QueryStmt{Select: &Select{Table: "test", Limit: &Limit{Count: 100}}},
		},
		{
			input:    `r test offset 100`,
			wantStmt: &QueryStmt{Select: &Select{Table: "test", Offset: &Offset{Count: 100}}},
		},
		{
			input:    `r test limit 200 offset 100`,
			wantStmt: &QueryStmt{Select: &Select{Table: "test", Limit: &Limit{Count: 200}, Offset: &Offset{Count: 100}}},
		},
		{
			input:    `r test offset 100 limit 200`,
			wantStmt: &QueryStmt{Select: &Select{Table: "test", Limit: &Limit{Count: 200}, Offset: &Offset{Count: 100}}},
		},
		{
			input:   `r test limit a`,
			wantErr: true,
		},
		{
			input:   `r test offset a`,
			wantErr: true,
		},
		{
			input:   `r test limit a offset 5`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		q, err := parse(tc.input)
		if tc.wantErr != (err != nil) {
			t.Fatalf("wantErr: %v, err: %v", tc.wantErr, err)
		}

		if !reflect.DeepEqual(tc.wantStmt, q) {
			t.Fatalf("want: %+v, got: %+v", tc.wantStmt, q)
		}
	}
}

func TestParse_Insert(t *testing.T) {
	tests := []struct {
		input    string
		wantStmt *QueryStmt
		wantErr  bool
	}{
		{
			input:   "insert test",
			wantErr: true,
		},
		{
			input:   "insert into test",
			wantErr: true,
		},
		{
			input:   "insert into test (aaa, bbb)",
			wantErr: true,
		},
		{
			input:   "insert into test (aaa, bbb) values",
			wantErr: true,
		},
		{
			input:   "insert into test (aaa, bbb) values 'ccc', 'ddd'",
			wantErr: true,
		},
		{
			input:   "insert into test (aaa, bbb) values 'ccc' 'ddd'",
			wantErr: true,
		},
		{
			input:   "insert into test (aaa, bbb) values ('ccc' 'ddd'",
			wantErr: true,
		},
		{
			input:   "insert into test (aaa, bbb values ('ccc' 'ddd')",
			wantErr: true,
		},
		{
			input:    "insert into test (aaa, bbb) values ('ccc', 'ddd')",
			wantStmt: &QueryStmt{Insert: &Insert{Table: "test", Cols: []string{"aaa", "bbb"}, Vals: []string{"ccc", "ddd"}}},
		},
		{
			input:    `insert into test (aaa, bbb) values ("ccc", "ddd")`,
			wantStmt: &QueryStmt{Insert: &Insert{Table: "test", Cols: []string{"aaa", "bbb"}, Vals: []string{"ccc", "ddd"}}},
		},
		{
			input:    `insert into test (aaa, bbb) values ("ccc", "ddd")`,
			wantStmt: &QueryStmt{Insert: &Insert{Table: "test", Cols: []string{"aaa", "bbb"}, Vals: []string{"ccc", "ddd"}}},
		},
		{
			input:    `insert into test values ("ccc", "ddd")`,
			wantStmt: &QueryStmt{Insert: &Insert{Table: "test", Cols: []string{}, Vals: []string{"ccc", "ddd"}}},
		},
	}

	for _, tc := range tests {
		q, err := parse(tc.input)
		if tc.wantErr != (err != nil) {
			t.Fatalf("wantErr: %v, err: %v", tc.wantErr, err)
		}

		if !reflect.DeepEqual(tc.wantStmt, q) {
			t.Fatalf("want: %+v, got: %+v", tc.wantStmt, q)
		}
	}
}

func TestParse_Create(t *testing.T) {
	tests := []struct {
		input    string
		wantStmt *QueryStmt
		wantErr  bool
	}{
		{
			input:   "create",
			wantErr: true,
		},
		{
			input:   "create table",
			wantErr: true,
		},
		{
			input:   "create table test",
			wantErr: true,
		},
		{
			input:   "create table test col1 string",
			wantErr: true,
		},
		{
			input:    "create table test (col1 string)",
			wantStmt: &QueryStmt{Create: &Create{Table: "test", Cols: []string{"col1"}, Types: []string{"string"}}},
		},
		{
			input:    "create table test ( col1 string , col2 string )",
			wantStmt: &QueryStmt{Create: &Create{Table: "test", Cols: []string{"col1", "col2"}, Types: []string{"string", "string"}}},
		},
		{
			input:    "create table test (col1 string,col2 string)",
			wantStmt: &QueryStmt{Create: &Create{Table: "test", Cols: []string{"col1", "col2"}, Types: []string{"string", "string"}}},
		},
		{
			input:    "create table test (col1 string ,col2 string)",
			wantStmt: &QueryStmt{Create: &Create{Table: "test", Cols: []string{"col1", "col2"}, Types: []string{"string", "string"}}},
		},
		{
			input:   "create table test (col1 string ,col2 int)",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		q, err := parse(tc.input)
		if tc.wantErr != (err != nil) {
			t.Fatalf("wantErr: %v, err: %v", tc.wantErr, err)
		}

		if !reflect.DeepEqual(tc.wantStmt, q) {
			t.Fatalf("want: %+v, got: %+v", tc.wantStmt, q)
		}
	}
}
