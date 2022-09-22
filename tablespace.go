package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// schema
// {
//   "tbl1": [
//     {"key1": "val1"},
//     {"key2": "val2"},
//     {"key3": "val3"}
//   ],
//   "tbl2": [
//     {"key1": "val1"},
//     {"key2": "val2"},
//     {"key3": "val3"}
//   ]
// }
var datafile = "data/incdb.data"

func init() {
	if os.Getenv("INCDB_TEST") == "1" {
		datafile = "data/test.incdb.data"
	}
}

func readAll(tbl string) ([]map[string]string, error) {
	f, err := os.OpenFile(datafile, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("open tablespace file: %w", err)
	}
	defer f.Close()

	d := map[string][]map[string]string{}

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat tablespace file: %w", err)
	}

	if info.Size() == 0 {
		return nil, fmt.Errorf("table '%s' not found", tbl)
	}

	// Do the JSON decode in case the file is not empty. It will be empty only on the incdb first time run.
	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return nil, fmt.Errorf("decode tablespace file as JSON: %w", err)
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

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat tablespace file: %w", err)
	}

	if info.Size() == 0 {
		return fmt.Errorf("table '%s' not found", tbl)
	}

	// Do the JSON decode in case the file is not empty. It will be empty only on the incdb first time run.
	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return fmt.Errorf("decode tablespace file as JSON: %w", err)
	}

	if _, ok := d[tbl]; !ok {
		return fmt.Errorf("table '%s' not found", tbl)
	}

	d[tbl] = append(d[tbl], map[string]string{key: value})

	// Drop the file content before write. Seek(0, 0) is needed to modify the IO offset.
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("clear tablespace file: %w", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("change tablespace file IO offset: %w", err)
	}

	// Write the data back.
	if err := json.NewEncoder(f).Encode(d); err != nil {
		return fmt.Errorf("encode data into tablespace file: %w", err)
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

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat tablespace file: %w", err)
	}

	// Do the JSON decode in case the file is not empty. It will be empty only on the incdb first time run.
	if info.Size() != 0 {
		if err := json.NewDecoder(f).Decode(&d); err != nil {
			return fmt.Errorf("decode tablespace file as JSON: %w", err)
		}
	}

	if _, ok := d[tbl]; ok {
		return fmt.Errorf("table '%s' already exists", tbl)
	}

	d[tbl] = []map[string]string{}

	// Drop the file content before write. Seek(0, 0) is needed to modify the IO offset.
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("clear tablespace file: %w", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("change tablespace file IO offset: %w", err)
	}

	// Write the data back.
	if err := json.NewEncoder(f).Encode(d); err != nil {
		return fmt.Errorf("encode data into tablespace file: %w", err)
	}

	f.Sync()
	return nil
}
