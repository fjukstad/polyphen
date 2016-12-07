package main

import (
	"flag"
	"fmt"
)

type Variant struct {
	Chromosome  string
	Position    string
	Id          string
	Reference   string
	Alternative string
	Qual        string
	Filter      string
	Info        string
	Format      string
	Additional  []string
}

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

	variants, err := parseVcf(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(variants)
	u := "http://genetics.bwh.harvard.edu/cgi-bin/ggi/ggi2.cgi"
	fmt.Println(u)
}
