package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func TestE2E(t *testing.T) {
	tests := []struct {
		query string
		// used in insert/create
		msg string
		// used in select
		rHdr []string
		rDat [][]string
		// file verification
		data map[string][]map[string]string
		cat  *Catalog
	}{
		// create table
		{
			query: "create table item (id string, name string)",
			msg:   "table item created",
			data:  map[string][]map[string]string{"item": {}},
			cat: &Catalog{Tables: []*CtTable{
				{
					Name: "item",
					Cols: []*CtCol{
						{Name: "id", Type: "string"},
						{Name: "name", Type: "string"},
					},
				},
			}},
		},
		{
			query: "create table user (id string, name string, city string)",
			msg:   "table user created",
			data:  map[string][]map[string]string{"item": {}, "user": {}},
			cat: &Catalog{Tables: []*CtTable{
				{
					Name: "item",
					Cols: []*CtCol{
						{Name: "id", Type: "string"},
						{Name: "name", Type: "string"},
					},
				},
				{
					Name: "user",
					Cols: []*CtCol{
						{Name: "id", Type: "string"},
						{Name: "name", Type: "string"},
						{Name: "city", Type: "string"},
					},
				},
			}},
		},

		// insert
		{
			query: `insert into item (id, name) values ("1", "laptop")`,
			msg:   "inserted",
			data: map[string][]map[string]string{"user": {}, "item": {
				{"id": "1", "name": "laptop"},
			}},
		},
		{
			query: `insert into item (id, name) values ("2", "iPhone")`,
			msg:   "inserted",
			data: map[string][]map[string]string{"user": {}, "item": {
				{"id": "1", "name": "laptop"},
				{"id": "2", "name": "iPhone"},
			}},
		},
		{
			query: `insert into item values ("3", "radio")`,
			msg:   "inserted",
			data: map[string][]map[string]string{"user": {}, "item": {
				{"id": "1", "name": "laptop"},
				{"id": "2", "name": "iPhone"},
				{"id": "3", "name": "radio"},
			}},
		},
		{
			query: `insert into item (id) values ("4")`,
			msg:   "inserted",
			data: map[string][]map[string]string{"user": {}, "item": {
				{"id": "1", "name": "laptop"},
				{"id": "2", "name": "iPhone"},
				{"id": "3", "name": "radio"},
				{"id": "4"},
			}},
		},

		// select
		{
			query: "r item",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"1", "laptop"},
				{"2", "iPhone"},
				{"3", "radio"},
				{"4", ""},
			},
		},
		{
			query: "r item '1'",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"1", "laptop"},
			},
		},
		{
			query: "r item '99'",
			msg:   "no results",
		},
		{
			query: "r item limit 2",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"1", "laptop"},
				{"2", "iPhone"},
			},
		},
		{
			query: "r item limit 100",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"1", "laptop"},
				{"2", "iPhone"},
				{"3", "radio"},
				{"4", ""},
			},
		},
		{
			query: "r item offset 1",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"2", "iPhone"},
				{"3", "radio"},
				{"4", ""},
			},
		},
		{
			query: "r item offset 4",
			msg:   "no results",
		},
		{
			query: "r item limit 1 offset 1",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"2", "iPhone"},
			},
		},
		{
			query: "r item limit 3 offset 1",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"2", "iPhone"},
				{"3", "radio"},
				{"4", ""},
			},
		},
		{
			query: "r item offset 1 limit 3",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"2", "iPhone"},
				{"3", "radio"},
				{"4", ""},
			},
		},
		{
			query: "r item order by asc",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"1", "laptop"},
				{"2", "iPhone"},
				{"3", "radio"},
				{"4", ""},
			},
		},
		{
			query: "r item order by desc",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"4", ""},
				{"3", "radio"},
				{"2", "iPhone"},
				{"1", "laptop"},
			},
		},
		{
			query: "r item order by desc limit 2",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"4", ""},
				{"3", "radio"},
			},
		},
		{
			query: "r item order by desc limit 2 offset 1",
			rHdr:  []string{"id", "name"},
			rDat: [][]string{
				{"3", "radio"},
				{"2", "iPhone"},
			},
		},
	}

	// prepare test
	os.Setenv("INCDB_TEST", "1")
	t.Cleanup(func() { os.Unsetenv("INCDB_TEST") })

	exec.Command("rm", "-f", "./data/test.incdb.data").Run()
	exec.Command("rm", "-f", "./data/test.incdb.catalog").Run()
	incdbd := exec.Command("./incdbd")
	if err := incdbd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		process, err := os.FindProcess(incdbd.Process.Pid)
		if err != nil {
			t.Fatal(err)
		}
		if err = process.Kill(); err != nil {
			t.Fatal(err)
		}
	})

	for _, tc := range tests {
		out, err := exec.Command("./incdb", tc.query).CombinedOutput()
		if err != nil {
			t.Fatalf("[%s] err: out: %s", tc.query, string(out))
		}

		// verify output
		if tc.msg != "" {
			// check message output
			if o := strings.TrimSuffix(string(out), "\n"); tc.msg != o {
				t.Fatalf("[%s] out: expected: '%s', got: '%s'", tc.query, tc.msg, o)
			}
		} else {
			// check query result
			type Output struct {
				Hdr  []string
				Vals [][]string
			}

			var o Output
			if err := json.Unmarshal(out, &o); err != nil {
				t.Fatalf("[%s] unmarshal output into json: %v", tc.query, err)
			}

			if !reflect.DeepEqual(tc.rHdr, o.Hdr) {
				t.Fatalf("[%s] out/header: expected: '%s', got: '%s'", tc.query, tc.rHdr, o.Hdr)
			}

			if !reflect.DeepEqual(tc.rDat, o.Vals) {
				t.Fatalf("[%s] out/data: expected: '%s', got: '%s'", tc.query, tc.rDat, o.Vals)
			}
		}

		// verify data file
		if tc.data != nil {
			d, err := os.ReadFile("./data/test.incdb.data")
			if err != nil {
				t.Fatalf("[%s] read test data file: %v", tc.query, err)
			}
			d = d[:len(d)-1] // trim new line

			e, err := json.Marshal(&tc.data)
			if err != nil {
				t.Fatalf("[%s] marshal expected data into json: %v", tc.query, err)
			}

			if !reflect.DeepEqual(d, e) {
				t.Fatalf("[%s] data: expected: '%s', got: '%s'", tc.query, string(e), string(d))
			}
		}

		// verify catalog file
		if tc.cat != nil {
			c, err := os.ReadFile("./data/test.incdb.catalog")
			if err != nil {
				t.Fatalf("[%s] read test catalog file: %v", tc.query, err)
			}
			c = c[:len(c)-1] // trim new line

			e, err := json.Marshal(&tc.cat)
			if err != nil {
				t.Fatalf("[%s] marshal expected catalog into json: %v", tc.query, err)
			}

			if !reflect.DeepEqual(c, e) {
				t.Fatalf("[%s] catalog: expected: '%s', got: '%s'", tc.query, string(e), string(c))
			}
		}
	}
}
