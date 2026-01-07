package utils

import (
	"encoding/json"
	"io"
)

type JsonObject = map[string]interface{}

func ReadJsonLines[E any](reader io.Reader) ([]E, error) {
	var items []E
	decoder := json.NewDecoder(reader)
	for decoder.More() {
		var item E
		if err := decoder.Decode(&item); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		} else {
			items = append(items, item)
		}
	}
	return items, nil
}

func WriteJsonLines[E any](writer io.Writer, items []E) error {
	encoder := json.NewEncoder(writer)
	for _, item := range items {
		if err := encoder.Encode(item); err != nil {
			return err
		}
	}
	return nil
}
