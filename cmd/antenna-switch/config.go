package main

import (
	"encoding/json"
	"os"
)

type Port struct {
	Label    string
	Position int
}
type Config struct {
	ListenAddress string
	Ports         []Port
}

func (c Config) WriteTo(destFile string) error {
	f, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}

func DefaultConfig() Config {
	return Config{
		ListenAddress: "0.0.0.0:8123",
		Ports: []Port{
			{"40m", 32},
			{"20m", 16},
			{"Ground", 0},
			{"40/20m", -16},
			{"6m", -32},
		},
	}
}
