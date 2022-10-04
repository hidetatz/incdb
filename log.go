package main

import (
	"fmt"
	"os"
)

var debug = true

func init() {
	// test mode
	if os.Getenv("INCDB_TEST") == "1" {
		debug = false
	}
}

func Debug(args ...any) {
	if debug {
		fmt.Fprint(os.Stdout, "[DEBUG] ")
		fmt.Fprintln(os.Stdout, args...)
	}
}
