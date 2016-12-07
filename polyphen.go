package main

import (
	"encoding/csv"
	"io/ioutil"
	"net/http"
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

	header := []string{"# Chromosome:position", "Reference/Variant nucleotides"}
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

func getStatusMessage(id string) string {
	endpoint := "http://genetics.bwh.harvard.edu/ggi/pph2/"
	baseUrl := endpoint + id + "/1/"
	resp, err := http.Get(baseUrl + "started.txt")
	if err != nil {
		return "Batch not started. Check back later"
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	statusMessage := string(responseBody)
	resp, err = http.Get(baseUrl + "completed.txt")
	if err == nil {
		responseBody, err = ioutil.ReadAll(resp.Body)
		statusMessage = statusMessage + string(responseBody)
	}
	return statusMessage
}
