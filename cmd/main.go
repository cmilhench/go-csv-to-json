package main

import (
	"context"
	"log"
	"os"

	"github.com/cm-dev/template/internal/core/csv"
	"github.com/cm-dev/template/internal/core/transformer"
	"github.com/cm-dev/template/internal/core/writer"
	"github.com/cm-dev/template/internal/ports"
)

func main() {
	ctx := context.Background()
	// Extract
	input, err := os.Open("input.csv")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer input.Close()

	// var t ports.Transformer[any] = &transformer.JsonTransformer{}
	// var w ports.Writer[any] = &writer.JsonWriter{}

	var t ports.Transformer[[]string] = &transformer.CsvTransformer{}
	var w ports.Writer[[]string] = &writer.CsvWriter{}

	// Transform
	lines, err1 := csv.Read(ctx, input)
	items, err2 := t.Transform(ctx, lines)

	go func() {
		for err := range err1 {
			log.Panicln(err)
		}
	}()
	go func() {
		for err := range err2 {
			log.Panicln(err)
		}
	}()

	// Load
	err = w.Write(ctx, os.Stdout, items)
	for err != nil {
		log.Panicln(err)
	}
}
