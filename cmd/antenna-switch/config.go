package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	ListenAddress string
	Ports         []string
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
		Ports:         []string{"40M", "20M", "Ground", "Unused", "Unused"},
	}
}
