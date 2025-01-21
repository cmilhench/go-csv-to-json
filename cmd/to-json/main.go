package main

import (
	"context"
	"log"
	"os"

	"github.com/cm-dev/template/internal/adapter/opener"
	"github.com/cm-dev/template/internal/core/json"
	"github.com/cm-dev/template/internal/core/pipeline"
	"github.com/cm-dev/template/internal/core/writer"
)

var head []string

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
	items := pipeline.Map(ctx, lines, Prettify(<-pipeline.Take(ctx, lines, 1)))

	go func() {
		for err := range errs {
			log.Panicln(err)
		}
	}()

	// Load
	err = writer.WriteJson(ctx, os.Stdout, items)
	for err != nil {
		log.Panicln(err)
	}
}

func Prettify(head []string) func(context.Context, int, []string) any {
	return func(ctx context.Context, i int, line []string) any {
		out, _ := json.Inflate(head, line)
		return out
	}
}
