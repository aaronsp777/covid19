#!/bin/sh
cd COVID-19
git pull
cd -
./builder > static/global.csv
./builder --incremental > static/global_rates.csv
./builder -in COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_US.csv -region_index 6 -counts_index 11  > static/us.csv
./builder -in COVID-19/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_US.csv -region_index 6 -counts_index 11 --incremental > static/us_rates.csv
