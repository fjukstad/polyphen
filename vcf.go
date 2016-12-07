package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
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

func parseVcf(filename string) ([]Variant, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []Variant{}, errors.Wrap(err, "Error: Can't open file"+filename)
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
		if len(fields) < 2 {
			return []Variant{}, errors.New("Error: could not parse vcf file. Something is wrong here:" + line)
		}
		// grab the first part of each line (contains chromosome, ref, alt)
		variantInfo := strings.Split(fields[0], "\t")
		format := fields[1]
		var additionalFields []string
		if len(fields) > 2 {
			additionalFields = fields[2 : len(fields)-1]
		}

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
	return variants, nil
}
