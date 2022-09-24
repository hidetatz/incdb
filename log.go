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

func Debug(arg any) {
	if debug {
		fmt.Fprintf(os.Stdout, "[DEBUG] %+v\n", arg)
	}
}
