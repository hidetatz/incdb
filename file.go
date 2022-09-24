package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func readJsonFile(f *os.File, dst any) error {
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat tablespace file: %w", err)
	}

	if info.Size() == 0 {
		return nil
	}

	if err := json.NewDecoder(f).Decode(dst); err != nil {
		return fmt.Errorf("decode tablespace file as JSON: %w", err)
	}

	return nil
}

func updateJsonFile(f *os.File, src any) error {
	// Drop the file content before write. Seek(0, 0) is needed to modify the IO offset.
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("clear tablespace file: %w", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("change tablespace file IO offset: %w", err)
	}

	if err := json.NewEncoder(f).Encode(src); err != nil {
		return fmt.Errorf("encode data into file: %w", err)
	}

	return nil
}
