package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func read(key string) (string, error) {
	d, err := decode()
	if err != nil {
		return "", fmt.Errorf("decode file: %w", err)
	}

	v, ok := d[key]
	if !ok {
		return "not found", nil
	}

	return fmt.Sprintf("%s", v), nil
}

func readAll() (string, error) {
	d, err := decode()
	if err != nil {
		return "", fmt.Errorf("decode file: %w", err)
	}

	return fmt.Sprintf("%+v", d), nil
}

func decode() (map[string]string, error) {
	f, err := os.OpenFile("data/incdb.data", os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("open tablespace file: %w", err)
	}
	defer f.Close()

	d := make(map[string]string)

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat tablespace file: %w", err)
	}

	// Do the JSON decode in case the file is not empty. It will be empty only on the incdb first time run.
	if info.Size() == 0 {
		return d, nil
	}

	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return nil, fmt.Errorf("decode tablespace file as JSON: %w", err)
	}

	return d, nil
}

func save(key, value string) error {
	// Open tablespace file to read. The data format is JSON {key: value}
	f, err := os.OpenFile("data/incdb.data", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("open tablespace file: %w", err)
	}
	defer f.Close()

	d := make(map[string]string)

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

	// key is ok to be replaced if duplicate.
	d[key] = value

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
