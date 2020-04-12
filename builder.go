package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	inFile = flag.String("in",
		"COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv",
		"file to read")
	regionsString = flag.String("regions",
		"China,Australia,Canada,United Kingdom",
		"Regions to graph")
)

func main() {
	flag.Parse()

	regions := strings.Split(*regionsString, ",")

	f, err := os.Open(*inFile)
	if err == io.EOF {
		log.Fatal(err)
	}
	r := csv.NewReader(f)

	headers, err := r.Read()
	if err == io.EOF {
		log.Fatal(err)
	}
	dates := headers[4:]

	counts := make(map[string]int)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if len(record) < 4 {
			log.Print("error: short record:", record)
			continue
		}
		region := record[1]
		for i, c := range record[4:] {
			n, err := strconv.Atoi(c)
			if err != nil {
				log.Print(err)
			}
			k := region + "," + dates[i]
			counts[k] += n
		}
	}

	fmt.Print("Date")
	for _, r := range regions {
		fmt.Print(",", r)
	}
	fmt.Print("\n")

	for _, d := range dates {
		fmt.Print(d)
		for _, r := range regions {
			fmt.Print(",", counts[r+","+d])
		}
		fmt.Print("\n")
	}
}
