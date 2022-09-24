package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestE2E(t *testing.T) {
	tests := []struct {
		input        string
		wantErr      bool
		expectedOut  string
		expectedData string
	}{
		{input: "r test", expectedOut: "compute result: table 'test' not found", wantErr: true},
		{input: "create table test (key string, value string)", expectedData: `{"test":[]}`},
		{input: "insert into test", expectedOut: "gramatically invalid: parse statement: values is expected but got EOF", wantErr: true},
		{input: `insert into test (key, value) values ("1", "a")`, expectedData: `{"test":[{"key":"1","value":"a"}]}`},
		{input: `insert into test (key, value) values ("2", "b")`, expectedData: `{"test":[{"key":"1","value":"a"},{"key":"2","value":"b"}]}`},
		{input: `insert into test (key, value) values ("3", "c")`, expectedData: `{"test":[{"key":"1","value":"a"},{"key":"2","value":"b"},{"key":"3","value":"c"}]}`},
		{input: "r test", expectedOut: `[1 a][2 b][3 c]`},
		{input: "r test '1'", expectedOut: `[1 a]`},
		{input: "r test '99'", expectedOut: ``},
		{input: "r test limit 2", expectedOut: `[1 a][2 b]`},
		{input: "r test limit 100", expectedOut: `[1 a][2 b][3 c]`},
		{input: "r test limit", expectedOut: "gramatically invalid: parse statement: integer value is expected but got EOF", wantErr: true},
		{input: "r test limit a", expectedOut: "gramatically invalid: parse statement: integer value is expected but got symbol", wantErr: true},
		{input: "r test offset 1", expectedOut: `[2 b][3 c]`},
		{input: "r test offset 3", expectedOut: ""},
		{input: "r test offset", expectedOut: "gramatically invalid: parse statement: integer value is expected but got EOF", wantErr: true},
		{input: "r test offset abc", expectedOut: "gramatically invalid: parse statement: integer value is expected but got symbol", wantErr: true},
		{input: "r test limit 1 offset 1", expectedOut: "[2 b]"},
		{input: "r test limit 5 offset 1", expectedOut: "[2 b][3 c]"},
		{input: "r test offset 1 limit 5", expectedOut: "[2 b][3 c]"},
		{input: "r test order by asc", expectedOut: "[1 a][2 b][3 c]"},
		{input: "r test order by desc", expectedOut: "[3 c][2 b][1 a]"},
		{input: "r test order by desc limit 1", expectedOut: "[3 c]"},
		{input: "r test order by desc limit 1 offset 2", expectedOut: "[1 a]"},
		{
			input:        "create table test2 (key2 string, value2 string)",
			expectedData: `{"test":[{"key":"1","value":"a"},{"key":"2","value":"b"},{"key":"3","value":"c"}],"test2":[]}`,
		}, // flaky: the order of "test" and "test2" might be interchanged
		{
			input:        `insert into test2 (key2, value2) values ("1", "a")`,
			expectedData: `{"test":[{"key":"1","value":"a"},{"key":"2","value":"b"},{"key":"3","value":"c"}],"test2":[{"key2":"1","value2":"a"}]}`,
		},
		{
			input:        `insert into test2 (key2, value2) values ("2", "b")`,
			expectedData: `{"test":[{"key":"1","value":"a"},{"key":"2","value":"b"},{"key":"3","value":"c"}],"test2":[{"key2":"1","value2":"a"},{"key2":"2","value2":"b"}]}`,
		},
		{
			input:        `insert into test2 (key2, value2) values ("3", "c")`,
			expectedData: `{"test":[{"key":"1","value":"a"},{"key":"2","value":"b"},{"key":"3","value":"c"}],"test2":[{"key2":"1","value2":"a"},{"key2":"2","value2":"b"},{"key2":"3","value2":"c"}]}`,
		},
		{
			input:        `insert into test2 (key2) values ("4")`,
			expectedData: `{"test":[{"key":"1","value":"a"},{"key":"2","value":"b"},{"key":"3","value":"c"}],"test2":[{"key2":"1","value2":"a"},{"key2":"2","value2":"b"},{"key2":"3","value2":"c"},{"key2":"4"}]}`,
		},
		{
			input:        `insert into test2 values ("5", "e")`,
			expectedData: `{"test":[{"key":"1","value":"a"},{"key":"2","value":"b"},{"key":"3","value":"c"}],"test2":[{"key2":"1","value2":"a"},{"key2":"2","value2":"b"},{"key2":"3","value2":"c"},{"key2":"4"},{"key2":"5","value2":"e"}]}`,
		},
		{input: "r test2", expectedOut: `[1 a][2 b][3 c][4 ][5 e]`},
	}

	// prepare test
	os.Setenv("INCDB_TEST", "1")
	t.Cleanup(func() {
		os.Unsetenv("INCDB_TEST")
	})

	o, err := exec.Command("rm", "-f", "./data/test.incdb.data").CombinedOutput()
	if err != nil {
		t.Fatalf("rm data: %v", string(o))
	}

	o, err = exec.Command("rm", "-f", "./data/test.incdb.catalog").CombinedOutput()
	if err != nil {
		t.Fatalf("rm data: %v", string(o))
	}

	// do check. Not like an ordinary table driven test, each test runs sequentially on the same tablespace file.
	for _, tc := range tests {
		out, err := exec.Command("./incdb", tc.input).CombinedOutput()
		if tc.wantErr != (err != nil) {
			t.Fatalf("running command '%v' fail: %v, out: %v", tc.input, err, string(out))
		}

		if tc.expectedOut != "" {
			if o := strings.TrimSuffix(string(out), "\n"); tc.expectedOut != o {
				t.Fatalf("out: expected: '%s', got: '%s' (input: %s)", tc.expectedOut, o, tc.input)
			}
		}

		if tc.expectedData != "" {
			d, err := os.ReadFile("./data/test.incdb.data")
			if err != nil {
				t.Fatalf("read test data file: %v", err)
			}

			if dt := strings.TrimSuffix(string(d), "\n"); tc.expectedData != dt {
				t.Fatalf("data: expected: '%s', got: '%s' (input: %s)", tc.expectedData, dt, tc.input)
			}
		}

	}
}
