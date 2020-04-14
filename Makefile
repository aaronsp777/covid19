
default: server dygraph builder counts
pull:
	cd COVID-19 && git pull
clean:
	rm -f static/*.csv builder server

builder: builder.go
	go build $<
server: server.go
	go build $<

counts: static/global.csv static/global_rates.csv static/us.csv static/us_rates.csv
dygraph: static/dygraph.min.js static/dygraph.js static/dygraph.css

static/dygraph.min.js:
	cd static && curl -O http://dygraphs.com/2.1.0/dygraph.min.js
static/dygraph.js:
	cd static && curl -O http://dygraphs.com/2.1.0/dygraph.js
static/dygraph.css:
	cd static && curl -O http://dygraphs.com/2.1.0/dygraph.css

static/global.csv: COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv
	./builder -in $< > $@
static/global_rates.csv: COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv
	./builder -in $< --incremental > $@
static/us.csv: COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_US.csv
	./builder -in $< -region_index 6 -counts_index 11 > $@
static/us_rates.csv: COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_US.csv
	./builder -in $< -region_index 6 -counts_index 11 --incremental > $@

