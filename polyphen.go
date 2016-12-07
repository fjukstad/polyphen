package main

import (
	"encoding/csv"
	"io"
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

var endpoint = "http://genetics.bwh.harvard.edu/ggi/pph2/"

func getStatusMessage(id string) string {
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

func downloadResults(id, outputDirectory string) error {

	err := os.MkdirAll(outputDirectory, 0764)
	if err != nil {
		return err
	}

	baseUrl := endpoint + id + "/1/"

	filenames := []string{
		"pph2-short.txt",
		"pph2-full.txt",
		"pph2-snps.txt",
		"pph2-log.txt"}

	for _, filename := range filenames {
		output, err := os.Create(outputDirectory + "/" + filename)
		if err != nil {
			return errors.Wrap(err, "Could not create output file")
		}
		defer output.Close()
		u := baseUrl + filename
		response, err := http.Get(u)
		if err != nil {
			return errors.Wrap(err, "Could not download output file")
		}
		defer response.Body.Close()

		_, err = io.Copy(output, response.Body)
		if err != nil {
			return errors.Wrap(err, "Could not read output file")
		}
	}
	return nil
}
