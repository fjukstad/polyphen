package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
)

func main() {

	filename := flag.String("f", "input.vcf", "input query file")
	formatted := flag.Bool("formatted", false, "is the input file already formatted to be used in polyphen2?")
	email := flag.String("email", "", "e-mail for use in polyphen notification service")
	modelName := flag.String("modelname", "HumDiv", "Classifier model")
	ucscDb := flag.String("UCSCDB", "hg19", "Genome assembly")
	snpFunc := flag.String("snpfunc", " ", "Annotations” option. Can be m for missense, c for coding, or empty for all")
	snpFilter := flag.String("snpfilter", "1", "Transcripts option. Can be 0 for all, 1 for canonical, or 3 for CCDS")
	status := flag.String("status", "", "Get status for Session ID")
	download := flag.String("download", "", "Download results for Session ID")
	outputDirectory := flag.String("o", "output", "output directory to store results")

	flag.Parse()

	if *status != "" {
		statusMessage := getStatusMessage(*status)
		fmt.Println("Batch query status:")
		fmt.Println(statusMessage)
		return
	}

	if *download != "" {
		err := downloadResults(*download, *outputDirectory)
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

	sessionId, err := submitQuery(batchQueryFilename, *modelName, *ucscDb,
		*email, *snpFunc, *snpFilter)
	if err != nil {
		fmt.Println(err)
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
