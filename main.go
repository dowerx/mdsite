package main

import (
	"flag"
	"io/fs"
	"mdsite/doc"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inDir := flag.String("in", "src", "input directory")
	outDir := flag.String("out", "dist", "output directory")

	flag.Parse()

	err := filepath.WalkDir(*inDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() && strings.EqualFold(filepath.Ext(path), ".md") {
			outFile := filepath.Join(
				*outDir,
				strings.TrimSuffix(
					strings.TrimPrefix(path, *inDir),
					filepath.Ext(path))+
					".html")

			if err := os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
				return err
			}

			d, err := doc.Load(path)
			if err != nil {
				return err
			}

			html, err := d.HTML()
			if err != nil {
				return err
			}

			return os.WriteFile(outFile, html, 0644)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
