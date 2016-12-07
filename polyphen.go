package main

import (
	"encoding/csv"
	"os"

	"github.com/pkg/errors"
)

func writeBatchQuery(variants []Variant, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "Could not open file to write batch query")
	}
	w := csv.NewWriter(f)
	w.Comma = '\t'

	header := []string{"Chromosome:position", "Reference/Variant nucleotides"}
	err = w.Write(header)
	if err != nil {
		return errors.Wrap(err, "Could not write record to batch query file")
	}
	for _, variant := range variants {
		record := []string{"chr" + variant.Chromosome + ":" + variant.Position,
			variant.Reference + "/" + variant.Alternative}
		err = w.Write(record)
		if err != nil {
			return errors.Wrap(err, "Could not write record to batch query file")
		}
	}
	w.Flush()
	return w.Error()

}
