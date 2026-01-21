package main

import (
	"log"

	"github.com/Quak1/lr-exporter/internal/exporter"
)

func main() {
	cfg, err := exporter.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	if err := exporter.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
