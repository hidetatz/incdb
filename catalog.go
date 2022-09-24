package main

import (
	"fmt"
	"os"
)

// schema
// {
//   "tables": [
//     {
//       "name": "tbl1",
//       "cols": [
//         {"name": "col1", "type": "string"},
//         {"name": "col2", "type": "string"},
//         {"name": "col3", "type": "string"}
//       ]
//     },
//     {
//       "name": "tbl2",
//       "cols": [
//         {"name": "col4", "type": "string"},
//         {"name": "col5", "type": "string"},
//         {"name": "col6", "type": "string"},
//       ]
//     }
//   ]
// }

var catfile = "data/incdb.catalog"

func init() {
	// test mode
	if os.Getenv("INCDB_TEST") == "1" {
		datafile = "data/test.incdb.catalog"
	}

	// initialize file if empty
	f, err := os.OpenFile(catfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read catalog file for initialization: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read catalog file stat for initialization: %s\n", err)
		os.Exit(1)
	}

	if info.Size() != 0 {
		return
	}

	if _, err := f.Write([]byte("{}")); err != nil {
		fmt.Fprintf(os.Stderr, "could not initialize empty catalog file: %s\n", err)
		os.Exit(1)
	}
}

type CtCol struct {
	Name string
	Type string
}

type CtTable struct {
	Name string
	Cols []*CtCol
}

type Catalog struct {
	Tables []*CtTable
}

func readCatalog(tbl string) (*CtTable, error) {
	f, err := os.OpenFile(datafile, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("open tablespace file: %w", err)
	}
	defer f.Close()

	c := &Catalog{}

	if err := readJsonFile(f, c); err != nil {
		return nil, fmt.Errorf("read catalog file: %w", err)
	}

	for _, t := range c.Tables {
		if t.Name == tbl {
			return t, nil
		}
	}

	return nil, fmt.Errorf("table '%s' not found in catalog", tbl)
}
