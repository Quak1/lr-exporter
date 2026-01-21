package exporter

import "time"

type Config struct {
	CatalogPath     string
	DestinationPath string
	StartDate       time.Time
	EndDate         time.Time
	Pick            int
	Rating          int
	Copy            bool
}

type Image struct {
	id               int
	path             string
	filename         string
	format           string
	sidecarExtension string
}
