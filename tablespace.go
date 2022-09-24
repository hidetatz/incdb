package main

import (
	"fmt"
	"os"
)

// schema
// {
//   "tbl1": [
//     {"col1": "val1", "col2": "val2", "col3": "val3"},
//     {"col1": "val4", "col2": "val5", "col3": "val6"},
//     {"col1": "val1", "col2": "val9", "col3": "val9"},
//   ],
//   "tbl2": [
//     {"col4": "val1", "col5": "val2", "col6": "val3"},
//     {"col4": "val4", "col5": "val5", "col6": "val6"},
//     {"col4": "val1", "col5": "val9", "col6": "val9"},
//   ]
// }
var datafile = "data/incdb.data"

func init() {
	// test mode
	if os.Getenv("INCDB_TEST") == "1" {
		datafile = "data/test.incdb.data"
	}

	// initialize file if empty
	f, err := os.OpenFile(datafile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read data file for initialization: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read data file stat for initialization: %s\n", err)
		os.Exit(1)
	}

	if info.Size() != 0 {
		return
	}

	if _, err := f.Write([]byte("{}")); err != nil {
		fmt.Fprintf(os.Stderr, "could not initialize empty data file: %s\n", err)
		os.Exit(1)
	}
}

func readData(tbl string) ([]map[string]string, error) {
	f, err := os.OpenFile(datafile, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("open tablespace file: %w", err)
	}
	defer f.Close()

	d := map[string][]map[string]string{}

	if err := readJsonFile(f, &d); err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	t, ok := d[tbl]
	if !ok {
		return nil, fmt.Errorf("table '%s' not found", tbl)
	}

	return t, nil
}

func save(tbl, key, value string) error {
	f, err := os.OpenFile(datafile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("open tablespace file: %w", err)
	}
	defer f.Close()

	d := map[string][]map[string]string{}

	if err := readJsonFile(f, &d); err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if _, ok := d[tbl]; !ok {
		return fmt.Errorf("table '%s' not found", tbl)
	}

	d[tbl] = append(d[tbl], map[string]string{key: value})

	if err := updateJsonFile(f, &d); err != nil {
		return fmt.Errorf("update tablespace file: %w", err)
	}

	f.Sync()
	return nil
}

func createTable(tbl string) error {
	f, err := os.OpenFile(datafile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("open tablespace file: %w", err)
	}
	defer f.Close()

	d := map[string][]map[string]string{}

	if err := readJsonFile(f, &d); err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if _, ok := d[tbl]; ok {
		return fmt.Errorf("table '%s' already exists", tbl)
	}

	d[tbl] = []map[string]string{}

	if err := updateJsonFile(f, &d); err != nil {
		return fmt.Errorf("update tablespace file: %w", err)
	}

	f.Sync()
	return nil
}
