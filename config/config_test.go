package config

import "testing"

func TestSave(t *testing.T) {
	if err := Save(); err != nil {
		panic(err)
	}
}
