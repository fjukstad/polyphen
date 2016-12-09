package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	if strings.Contains(statusMessage, "Object not found!") {
		return "Batch query not found for id " + id + ". If you've just submitted it, check back later!"
	}

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

func submitQuery(batchQueryFilename, modelName, ucscDb, email, snpFunc,
	snpFilter string) (string, error) {
	params := map[string]string{
		"_ggi_project":         "PPHWeb2",
		"_ggi_origin":          "query",
		"_ggi_target_pipeline": "1",
		"MODELNAME":            modelName,
		"UCSCDB":               ucscDb,
		"NOTIFYME":             email,
		"SNPFUNC":              snpFunc,
		"SNPFILTER":            snpFilter,
	}

	bf, err := os.Open(batchQueryFilename)
	if err != nil {
		return "", errors.Wrap(err, "Could not open batch query file")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("_ggi_batch_file", filepath.Base(batchQueryFilename))
	if err != nil {
		return "", errors.Wrap(err, "Could not create form file")
	}

	_, err = io.Copy(part, bf)
	if err != nil {
		return "", errors.Wrap(err, "Read batch query file")
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	writer.Close()

	u := "http://genetics.bwh.harvard.edu/cgi-bin/ggi/ggi2.cgi"
	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		return "", errors.Wrap(err, "Could not create polyphen2 request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Could not perform batch query to polyphen")
	}

	cookies := resp.Cookies()
	sessionId, err := getSessionId(cookies)
	if err != nil {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.Wrap(err, "Could not response from polyphen")
		}
		return "", errors.Wrap(err, string(responseBody))
	}
	return sessionId, nil
}
