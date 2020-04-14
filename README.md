# covid19
Graphing covid19 timeseries

# Get dataset

```sh
git clone https://github.com/CSSEGISandData/COVID-19
```

# Download dygraph and put the following files into static:

```sh
pushd static
curl -O http://dygraphs.com/2.1.0/dygraph.min.js
curl -O http://dygraphs.com/2.1.0/dygraph.js
curl -O http://dygraphs.com/2.1.0/dygraph.css
popd
```

# Build everything

```sh
make
```

# Run the server
```sh
./server &
```

# Visit http://localhost:8080/static/global.html
