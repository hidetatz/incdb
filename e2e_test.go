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
		{input: "r test", expectedOut: "read data: table 'test' not found", wantErr: true},
		{input: "create table test", expectedData: `{"test":[]}`},
		{input: "w test", expectedOut: "gramatically invalid: parse statement: not a string", wantErr: true},
		{input: "w test 1 a", expectedData: `{"test":[{"1":"a"}]}`},
		{input: "w test 2 b", expectedData: `{"test":[{"1":"a"},{"2":"b"}]}`},
		{input: "w test 3 c", expectedData: `{"test":[{"1":"a"},{"2":"b"},{"3":"c"}]}`},
		{input: "r test", expectedOut: `[map[1:a] map[2:b] map[3:c]]`},
		{input: "r test 1", expectedOut: `[map[1:a]]`},
		{input: "r test 99", expectedOut: `[]`},
		{input: "r test limit 2", expectedOut: `[map[1:a] map[2:b]]`},
		{input: "r test limit 100", expectedOut: `[map[1:a] map[2:b] map[3:c]]`},
		{input: "r test limit", expectedOut: "gramatically invalid: parse statement: not a string", wantErr: true},
		{input: "r test limit a", expectedOut: "gramatically invalid: parse statement: cannot convert to int", wantErr: true},
		{input: "r test offset 1", expectedOut: `[map[2:b] map[3:c]]`},
		{input: "r test offset 3", expectedOut: "[]"},
		{input: "r test offset", expectedOut: "gramatically invalid: parse statement: not a string", wantErr: true},
		{input: "r test offset abc", expectedOut: "gramatically invalid: parse statement: cannot convert to int", wantErr: true},
		{input: "r test limit 1 offset 1", expectedOut: "[map[2:b]]"},
		{input: "r test limit 5 offset 1", expectedOut: "[map[2:b] map[3:c]]"},
		{input: "r test offset 1 limit 5", expectedOut: "[map[2:b] map[3:c]]"},
		{input: "r test order by asc", expectedOut: "[map[1:a] map[2:b] map[3:c]]"},
		{input: "r test order by desc", expectedOut: "[map[3:c] map[2:b] map[1:a]]"},
		{input: "r test order by desc limit 1", expectedOut: "[map[3:c]]"},
		{input: "r test order by desc limit 1 offset 2", expectedOut: "[map[1:a]]"},
		{input: "create table test2", expectedData: `{"test":[{"1":"a"},{"2":"b"},{"3":"c"}],"test2":[]}`}, // flaky: the order of "test" and "test2" might be interchanged
		{input: "w test2 1 a", expectedData: `{"test":[{"1":"a"},{"2":"b"},{"3":"c"}],"test2":[{"1":"a"}]}`},
		{input: "w test2 2 b", expectedData: `{"test":[{"1":"a"},{"2":"b"},{"3":"c"}],"test2":[{"1":"a"},{"2":"b"}]}`},
		{input: "w test2 3 c", expectedData: `{"test":[{"1":"a"},{"2":"b"},{"3":"c"}],"test2":[{"1":"a"},{"2":"b"},{"3":"c"}]}`},
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

	// do check. Not like an ordinary table driven test, each test runs sequentially on the same tablespace file.
	for _, tc := range tests {
		c := strings.Split(tc.input, " ")

		out, err := exec.Command("./incdb", c[0:]...).CombinedOutput()
		if tc.wantErr != (err != nil) {
			t.Fatalf("running command '%v' fail: %v", tc.input, err)
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
