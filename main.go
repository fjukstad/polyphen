package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
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

	f, err := os.Open(*filename)
	if err != nil {
		fmt.Println("Error: Can't open file", *filename, err)
		return
	}
	defer f.Close()

	var variants []Variant

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ";")
		if len(fields) < 3 {
			fmt.Println("Error: could not parse vcf file. Something is wrong here:", line)
			return
		}
		// grab the first part of each line (contains chromosome, ref, alt)
		variantInfo := strings.Split(fields[0], "\t")
		format := fields[1]
		additionalFields := fields[2 : len(fields)-1]

		chromosome := variantInfo[0]
		pos := variantInfo[1]
		id := variantInfo[2]
		ref := variantInfo[3]
		alt := variantInfo[4]
		qual := variantInfo[5]
		filter := variantInfo[6]
		info := variantInfo[7]

		variants = append(variants, Variant{
			chromosome,
			pos,
			id,
			ref,
			alt,
			qual,
			filter,
			info,
			format,
			additionalFields,
		})

	}

	fmt.Println(variants)
	u := "http://genetics.bwh.harvard.edu/cgi-bin/ggi/ggi2.cgi"
	fmt.Println(u)
}
