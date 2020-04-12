package main

import (
	"encoding/csv"
	"flag"
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
	regionsString = flag.String(
		"regions",
		"China,Australia,Canada,United Kingdom",
		"Regions to graph")
	countsIndex = flag.Int(
		"counts_index",
		4,
		"First column with counts (zero based.)")
	regionIndex = flag.Int(
		"region_index",
		1,
		"Column of Region (zero based.)")
	calculateIncrementals = flag.Bool(
		"incrementals",
		false,
		"Calculate daily incrementals over previous day")
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
	dates := headers[*countsIndex:]

	counts := make(map[string]int)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		region := record[*regionIndex]
		for i, c := range record[*countsIndex:] {
			n, err := strconv.Atoi(c)
			if err != nil {
				log.Print(err)
			}
			k := region + "," + dates[i]
			counts[k] += n
		}
	}

	if *calculateIncrementals {
		incrementals := make(map[string]int)

		for i, d := range dates {
			if i == 0 {
				continue
			}
			d0 := dates[i-1]
			for _, r := range regions {
				incrementals[r+","+d] = counts[r+","+d] - counts[r+","+d0]
			}
		}

		counts = incrementals
	}

	// Write out the records
	w := csv.NewWriter(os.Stdout)

	// Write header
	var record []string
	record = append(record, "Date")

	for _, r := range regions {
		record = append(record, r)
	}

	w.Write(record)
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	for _, d := range dates {
		record = append(record[:0], d)
		for _, r := range regions {
			record = append(record, strconv.Itoa(counts[r+","+d]))
		}
		w.Write(record)
		if err := w.Error(); err != nil {
			log.Fatalln("error writing csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
