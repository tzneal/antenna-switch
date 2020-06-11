package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	cfgPath := flag.String("config", "~/.antenna-switch.json", "path to the default configuration file")
	flag.Parse()

	if strings.HasPrefix(*cfgPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("unable to open config: %s", err)
		}
		*cfgPath = home + (*cfgPath)[1:]
	}

	log.Println("reading config from", *cfgPath)
	var config Config
	cfgFile, err := os.Open(*cfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			config = DefaultConfig()
			if err = config.WriteTo(*cfgPath); err != nil {
				log.Fatalf("error creating default config file: %s", err)
			}
		}
	} else {
		defer cfgFile.Close()
		dec := json.NewDecoder(cfgFile)
		if err = dec.Decode(&config); err != nil {
			log.Fatalf("error reading config file %s: %s", *cfgPath, err)
		}
	}

	server, err := NewServer(config.Ports)
	if err != nil {
		log.Fatalf("error creating server: %s", err)
	}
	http.HandleFunc("/", server.ServeIndex)
	http.HandleFunc("/switch", server.SwitchPorts)
	http.HandleFunc("/calibrate", server.Calibrate)
	log.Printf("%d ports configured %v", len(config.Ports), config.Ports)
	log.Printf("listening on %s", config.ListenAddress)
	log.Fatal(http.ListenAndServe(config.ListenAddress, nil))
}
