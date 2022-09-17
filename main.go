package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "at least one argument is required\n")
		os.Exit(1)
	}

	stmt, err := parse(strings.Join(os.Args[1:], " "))
	if err != nil {
		fmt.Fprintf(os.Stderr, "gramatically invalid: %s\n", err)
		os.Exit(1)
	}

	switch {
	case stmt.Select != nil:
		out, err := read()
		if err != nil {
			fmt.Fprintf(os.Stderr, "read data: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(out)

	case stmt.Insert != nil:
		if err := save(stmt.Insert.Key, stmt.Insert.Val); err != nil {
			fmt.Fprintf(os.Stderr, "save data %s: %s\n", stmt.Insert.Key, err)
			os.Exit(1)
		}
	}

}
