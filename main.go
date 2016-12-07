package main

import (
	"bytes"
	"errors"
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

	filename := flag.String("f", "input.vcf", "input query file")
	formatted := flag.Bool("formatted", false, "is the input file already formatted to be used in polyphen2?")
	email := flag.String("email", "", "e-mail for use in polyphen notification service")
	modelName := flag.String("modelname", "HumDiv", "Classifier model")
	ucscDb := flag.String("UCSCDB", "hg19", "Genome assembly")
	snpFunc := flag.String("snpfunc", "m", "Annotations” option. Can be m for missense, c for coding, or empty for all")
	snpFilter := flag.String("snpfilter", "1", "Transcripts option. Can be 0 for all, 1 for canonical, or 3 for CCDS")
	status := flag.Bool("status", false, "get output status of running polyphen query")
	id := flag.String("id", "", "Polyphen session ID")
	download := flag.Bool("download", false, "download results")
	outputDirectory := flag.String("o", "output", "output directory to store results")

	flag.Parse()

	if *status {
		if *id == "" {
			fmt.Println("Error: Please provide a valid Session ID")
			return
		}
		statusMessage := getStatusMessage(*id)
		fmt.Println("Batch query status:")
		fmt.Println(statusMessage)
		return
	}

	if *download {
		if *id == "" {
			fmt.Println("Error: Please provide a valid Session ID")
			return
		}
		err := downloadResults(*id, *outputDirectory)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Output saved in", *outputDirectory)
		return
	}

	if *email == "" {
		fmt.Println("Error: Please specify an e-mail address")
		return
	}

	batchQueryFilename := "polyphen2-batchquery.txt"

	// if input file has been formatted we can simply post it to polyphen, if
	// not we'll need to parse it and post the parsed content
	if *formatted {
		batchQueryFilename = *filename
	} else {
		variants, err := parseVcf(*filename)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = writeBatchQuery(variants, batchQueryFilename)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	params := map[string]string{
		"_ggi_project":         "PPHWeb2",
		"_ggi_origin":          "query",
		"_ggi_target_pipeline": "1",
		"MODELNAME":            *modelName,
		"UCSCDB":               *ucscDb,
		"NOTIFYME":             *email,
		"SNPFUNC":              *snpFunc,
		"SNPFILTER":            *snpFilter,
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

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not perfor batch query to polyphen2", err)
		return
	}

	cookies := resp.Cookies()
	sessionId, err := getSessionId(cookies)
	if err != nil {
		fmt.Println("Error: Bad response from Polyphen:")
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		fmt.Println(string(responseBody))
		return
	}

	fmt.Println("Polyphen batch query submission completed successfully. You'll get an e-mail at",
		*email, "when the query is completed. Until then you can check the progress with session ID",
		sessionId)
	return
}

func getSessionId(cookies []*http.Cookie) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == "polyphenweb2" {
			return cookie.Value, nil
		}
	}
	return "", errors.New("Session ID not found in cookies")
}
