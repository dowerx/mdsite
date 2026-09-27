package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"mdsite/doc"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	inDir := flag.String("in", "src", "input directory")
	outDir := flag.String("out", "dist", "output directory")
	headPath := flag.String("head", "", "head file")
	showHelp := flag.Bool("help", false, "show usage")

	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		flag.Usage()
		return nil
	}

	slog.Info("args", "in", *inDir, "out", *outDir, "head", *headPath)

	var head []byte
	if *headPath != "" {
		var err error
		if head, err = os.ReadFile(*headPath); err != nil {
			return err
		}
		slog.Info("loaded head", "path", *headPath)
	}

	doc.Init()

	return filepath.WalkDir(*inDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() ||
			strings.HasPrefix(filepath.Base(path), "_") ||
			!strings.EqualFold(filepath.Ext(path), ".md") {
			return nil
		}

		outFile := outPath(*inDir, *outDir, path)

		slog.Info("processing", "in", path, "out", outFile)

		if err := os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
			return err
		}

		d, err := doc.Load(path)
		if err != nil {
			return err
		}

		html, err := d.HTML(head)
		if err != nil {
			return err
		}

		return os.WriteFile(outFile, html, 0644)
	})
}

func usage() {
	fmt.Fprintf(
		flag.CommandLine.Output(),
		"Usage: %s [flags]\n\nRenders Markdown pages to HTML.\n\nFlags:\n",
		filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func outPath(inDir, outDir, path string) string {
	rel := strings.TrimPrefix(path, inDir)
	return filepath.Join(outDir, strings.TrimSuffix(rel, filepath.Ext(path))+".html")
}
