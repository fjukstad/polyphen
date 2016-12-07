package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	filename := flag.String("f", "input.vcf", "input vcf file")
	email := flag.String("email", "", "e-mail for use in polyphen notification service")
	modelName := flag.String("modelname", "HumDiv", "Classifier model")
	ucscDb := flag.String("UCSCDB", "hg19", "Genome assembly")
	snpFunc := flag.String("snpfunc", "m", "Annotations” option. Can be m for missense, c for coding, or empty for all")
	//snpFilter := flag.String("snpfilter", "1", "Transcripts option. Can be 0 for all, 1 for canonical, or 3 for CCDS")

	flag.Parse()

	if *email == "" {
		fmt.Println("Error: Please specify an e-mail address")
		return
	}

	variants, err := parseVcf(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	batchQueryFilename := "polyphen2-batchquery.txt"

	err = writeBatchQuery(variants, batchQueryFilename)
	if err != nil {
		fmt.Println(err)
		return
	}

	params := map[string]string{
		"_ggi_project":         "PPHWeb2",
		"_ggi_origin":          "query",
		"_ggi_target_pipeline": "1",
		"MODELNAME":            *modelName,
		"UCSCDB":               *ucscDb,
		"NOTIFYME":             *email,
		"SNPFUNC":              *snpFunc,
	}

	bf, err := os.Open(batchQueryFilename)
	if err != nil {
		fmt.Println("Could not open batch query file", err)
		return
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("_ggi_batch_file", filepath.Base(batchQueryFilename))
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = io.Copy(part, bf)
	if err != nil {
		fmt.Println(err)
		return
	}

	for key, val := range params {
		err = writer.WriteField(key, val)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	writer.Close()

	u := "http://genetics.bwh.harvard.edu/cgi-bin/ggi/ggi2.cgi"
	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		fmt.Println("Could not create polyphen2 request", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	fmt.Println(body)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not perfor batch query to polyphen2", err)
		return
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(responseBody))

}
