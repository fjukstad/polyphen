package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mendelics/vcf"
)

func main() {

	filename := flag.String("f", "input.vcf", "input vcf file")
	email := flag.String("email", "", "e-mail for use in polyphen notification service")
	//modelName := flag.String("modelname", "HumDiv", "Classifier model")
	//uscdDb := flag.String("UCSCDB", "hg19", "Genome assembly")
	//snpFunc := flag.String("snpfunc", "", "Annotations” option. Can be m for missense, c for coding, or empty for all")
	//snpFilter := flag.Int("snpfilter", 1, "Transcripts option. Can be 0 for all, 1 for canonical, or 3 for CCDS")

	flag.Parse()

	if *email == "" {
		fmt.Println("Error: Please specify an e-mail address")
		return
	}

	vcfFile, err := os.Open(*filename)
	if err != nil {
		fmt.Println("Error: Can't open file", *filename, err)
		return
	}
	defer vcfFile.Close()

	validVariants := make(chan *vcf.Variant, 2000)     // buffered channel for correctly parsed variants
	invalidVariants := make(chan vcf.InvalidLine, 100) // buffered channel for variants that fail to parse

	go func() {
		err := vcf.ToChannel(vcfFile, validVariants, invalidVariants)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		// consume invalid variants channel asynchronously
		for invalid := range invalidVariants {
			fmt.Println("failed to parse line", invalid.Line, "with error", invalid.Err)
		}
	}()

	for variant := range validVariants {
		fmt.Println(variant)
	}

	u := "http://genetics.bwh.harvard.edu/cgi-bin/ggi/ggi2.cgi"
	fmt.Println(u)
}
