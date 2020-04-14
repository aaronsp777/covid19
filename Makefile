
default: server dygraph builder counts
pull:
	cd COVID-19 && git pull
clean:
	rm -f static/*.csv builder server

builder: builder.go
	go build $<
server: server.go
	go build $<

counts: static/global.csv static/global_rates.csv static/us.csv static/us_rates.csv static/china.csv static/china_rates.csv static/washington.csv static/washington_rates.csv static/canada.csv static/canada_rates.csv static/california.csv static/california_rates.csv
dygraph: static/dygraph.min.js static/dygraph.js static/dygraph.css

static/dygraph.min.js:
	cd static && curl -O http://dygraphs.com/2.1.0/dygraph.min.js
static/dygraph.js:
	cd static && curl -O http://dygraphs.com/2.1.0/dygraph.js
static/dygraph.css:
	cd static && curl -O http://dygraphs.com/2.1.0/dygraph.css

CONFIRMED_GLOBAL = COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv
CONFIRMED_US = COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_US.csv

static/global.csv: $(CONFIRMED_GLOBAL)
	./builder -in $< > $@
static/global_rates.csv: $(CONFIRMED_GLOBAL)
	./builder -in $< --incremental > $@
static/us.csv: $(CONFIRMED_US)
	./builder -in $< -region_index 6 -counts_index 11 > $@
static/us_rates.csv: $(CONFIRMED_US)
	./builder -in $< -region_index 6 -counts_index 11 --incremental > $@

static/china.csv: $(CONFIRMED_GLOBAL)
	./builder -in $< --filter_value China -region_index 0 > $@
static/china_rates.csv: $(CONFIRMED_GLOBAL)
	./builder -in $< --filter_value China -region_index 0 --incremental > $@

static/canada.csv: $(CONFIRMED_GLOBAL)
	./builder -in $< --filter_value Canada -region_index 0 > $@
static/canada_rates.csv: $(CONFIRMED_GLOBAL)
	./builder -in $< --filter_value Canada -region_index 0 --incremental > $@

static/washington.csv: $(CONFIRMED_US)
	./builder -in $< -region_index 5 -counts_index 11 --filter_index 6 --filter_value Washington > $@
static/washington_rates.csv: $(CONFIRMED_US)
	./builder -in $< -region_index 5 -counts_index 11 --filter_index 6 --filter_value Washington --incremental > $@

static/california.csv: $(CONFIRMED_US)
	./builder -in $< -region_index 5 -counts_index 11 --filter_index 6 --filter_value California > $@
static/california_rates.csv: $(CONFIRMED_US)
	./builder -in $< -region_index 5 -counts_index 11 --filter_index 6 --filter_value California --incremental > $@
