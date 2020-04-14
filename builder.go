package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	inFile = flag.String("in",
		"COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv",
		"file to read")
	regionsString = flag.String(
		"regions",
		"",
		"Only show certain reasons (|) separated")
	countsIndex = flag.Int(
		"counts_index",
		4,
		"First column with counts (zero based.)")
	regionIndex = flag.Int(
		"region_index",
		1,
		"Column of Region (zero based.)")
	calculateIncrementals = flag.Bool(
		"incremental",
		false,
		"Calculate incremental daily rates over previous day")
	topN = flag.Int(
		"top",
		10,
		"Only show the top N graphs")
	filterIndex = flag.Int(
		"filter_index",
		1,
		"Column to filter (zero based.)")
	filterValue = flag.String(
		"filter_value",
		"",
		"Optional value to fitler against")
)

func topFilter(n int, ranks map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}

	if n > len(ranks) {
		n = len(ranks)
	}

	var ss []kv
	for k, v := range ranks {
		ss = append(ss, kv{k, v})
	}

	// Then sorting the slice by value, higher first.
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// Print the x top values
	var labels []string
	for _, kv := range ss[:n] {
		labels = append(labels, kv.Key)
	}
	return labels
}

func main() {
	flag.Parse()

	regions := strings.Split(*regionsString, "|")
	if *regionsString == "" {
		regions = nil
	}

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
	regionNames := make(map[string]bool)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if *filterValue != "" && *filterValue != record[*filterIndex] {
			continue
		}

		region := record[*regionIndex]
		regionNames[region] = true
		for i, c := range record[*countsIndex:] {
			n, err := strconv.Atoi(c)
			if err != nil {
				log.Print(err)
			}
			k := region + "," + dates[i]
			counts[k] += n
		}
	}

	if len(regions) == 0 {
		for k := range regionNames {
			regions = append(regions, k)
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

	if *topN > 0 {
		max := make(map[string]int)
		// Calculate max each for each region
		for _, d := range dates {
			for _, r := range regions {
				k := r + "," + d
				if max[r] < counts[k] {
					max[r] = counts[k]
				}
			}
		}
		regions = topFilter(*topN, max)
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
