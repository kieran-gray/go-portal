package utils

import (
	"encoding/json"
	"io"
	"os"
)

func ReadFile[T any](filename string, t T) (T, error) {
	file, err := os.Open(filename)
	if err != nil {
		return t, err
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	err = json.Unmarshal(bytes, &t)
	if err != nil {
		return t, err
	}

	return t, nil
}
