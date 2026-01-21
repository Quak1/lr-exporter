package exporter

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func Run(cfg *Config) error {
	dsn := "file:" + cfg.CatalogPath + "?mode=ro"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := queryImages(db, cfg)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		img, err := scanImage(rows)
		if err != nil {
			return err
		}

		if img.format == "RAW" && img.sidecarExtension == "" {
			log.Printf("%s doesn't have a sidecar image. Skipping.\n", img.filename)
			continue
		}

		src, dst := buildPaths(img, cfg)
		fmt.Printf("%d: '%s' -> '%s'\n", img.id, src, dst)

		if cfg.Copy {
			if err = copyFile(src, dst); err != nil {
				log.Println(err)
			}
		}
	}

	return rows.Err()
}

func scanImage(rows *sql.Rows) (*Image, error) {
	var img Image
	err := rows.Scan(&img.id, &img.path, &img.filename, &img.format, &img.sidecarExtension)
	img.path = filepath.FromSlash(img.path)
	return &img, err
}

func buildPaths(img *Image, cfg *Config) (src, dst string) {
	newFilename := replaceExtension(img.filename, img.sidecarExtension)
	src = filepath.Join(img.path, newFilename)
	dst = filepath.Join(cfg.DestinationPath, newFilename)
	return src, dst
}

func replaceExtension(path, ext string) string {
	idx := strings.LastIndex(path, ".")
	if idx == -1 {
		return path + "." + ext
	}
	return path[:idx] + "." + ext
}
