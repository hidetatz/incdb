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
		{input: "r", expectedOut: ""},
		{input: "w", expectedOut: "gramatically invalid: parse statement: not a string", wantErr: true},
		{input: "w 1 a", expectedData: `[{"1":"a"}]`},
		{input: "w 2 b", expectedData: `[{"1":"a"},{"2":"b"}]`},
		{input: "w 3 c", expectedData: `[{"1":"a"},{"2":"b"},{"3":"c"}]`},
		{input: "r", expectedOut: `[map[1:a] map[2:b] map[3:c]]`},
		{input: "r 1", expectedOut: `a`},
		{input: "r 99", expectedOut: `not found`},
		{input: "r limit 2", expectedOut: `[map[1:a] map[2:b]]`},
		{input: "r limit 100", expectedOut: `[map[1:a] map[2:b] map[3:c]]`},
		{input: "r limit", expectedOut: "gramatically invalid: parse statement: not a string", wantErr: true},
		{input: "r limit a", expectedOut: "gramatically invalid: parse statement: cannot convert to int", wantErr: true},
		{input: "r offset 1", expectedOut: `[map[2:b] map[3:c]]`},
		{input: "r offset 3", expectedOut: ""},
		{input: "r offset", expectedOut: "gramatically invalid: parse statement: not a string", wantErr: true},
		{input: "r offset abc", expectedOut: "gramatically invalid: parse statement: cannot convert to int", wantErr: true},
		{input: "r limit 1 offset 1", expectedOut: "[map[2:b]]"},
		{input: "r limit 5 offset 1", expectedOut: "[map[2:b] map[3:c]]"},
		{input: "r offset 1 limit 5", expectedOut: "[map[2:b] map[3:c]]"},
	}

	// prepare test
	os.Setenv("INCDB_TEST", "1")
	t.Cleanup(func() {
		os.Unsetenv("INCDB_TEST")
	})

	o, err := exec.Command("rm", "-f", "./incdb_test").CombinedOutput()
	if err != nil {
		t.Fatalf("rm bin: %v", string(o))
	}

	o, err = exec.Command("rm", "-f", "./data/test.incdb.data").CombinedOutput()
	if err != nil {
		t.Fatalf("rm data: %v", string(o))
	}

	o, err = exec.Command("go", "build", "-o", "incdb_test", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("build bin: %v", string(o))
	}

	// do check. Not like an ordinary table driven test, each test runs sequentially on the same tablespace file.
	for _, tc := range tests {
		c := strings.Split(tc.input, " ")

		out, err := exec.Command("./incdb_test", c[0:]...).CombinedOutput()
		if tc.wantErr != (err != nil) {
			t.Fatalf("running command '%v' fail: %v", tc.input, err)
		}

		if tc.expectedOut != "" {
			if tc.expectedOut+"\n" != string(out) {
				t.Fatalf("out: expected: '%s', got: '%s' (input: %s)", tc.expectedOut, string(out), tc.input)
			}
		}

		if tc.expectedData != "" {
			d, err := os.ReadFile("./data/test.incdb.data")
			if err != nil {
				t.Fatalf("read test data file: %v", err)
			}

			if tc.expectedData+"\n" != string(d) {
				t.Fatalf("data: expected: '%s', got: '%s' (input: %s)", tc.expectedData, string(d), tc.input)
			}
		}

	}
}
