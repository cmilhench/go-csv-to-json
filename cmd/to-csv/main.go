package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/cm-dev/template/internal/adapter/opener"
	"github.com/cm-dev/template/internal/core/pipeline"
	"github.com/cm-dev/template/internal/core/writer"
)

func main() {
	ctx := context.Background()

	// Extract
	fs := opener.OpenerFS{}
	input, err := fs.Open(ctx, ".data/report")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer input.Close()

	// Transform
	lines, errs := pipeline.ReadCsv(ctx, input)
	lines = pipeline.Take(ctx, lines, 3)
	items := pipeline.Map(ctx, lines, Prettify)

	go func() {
		for err := range errs {
			log.Panicln(err)
		}
	}()

	// Load
	err = writer.WriteCsv(ctx, os.Stdout, items)
	for err != nil {
		log.Panicln(err)
	}
}

func Prettify(ctx context.Context, i int, line []string) []string {
	if i > 0 {
		return line
	}
	for i, v := range line {
		v = strings.ToUpper(string(v[0])) + v[1:] // capitalise
		v = strings.ReplaceAll(v, ".", " ")       // replace "."s with " "s
		line[i] = v
	}
	return line
}
