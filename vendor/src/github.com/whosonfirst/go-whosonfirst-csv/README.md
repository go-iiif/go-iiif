# go-whosonfirst-csv

A simple Go package to implement a "dict reader" style CSV parser (on top of the default `encoding/csv` package) to return rows a key-value dictionaries rather than lists.

## Setup

Run the handy `build` target in the [Makefile](Makefile).

```
make build
```

## Usage

### Reading files

```
import (
	csv "github.com/whosonfirst/go-whosonfirst-csv"
	"os"
)

reader, reader_err := csv.NewDictReaderFromPath("example.csv")

// or maybe you might do
// reader, err := csv.NewDictReader(os.Stdin)

if err != nil {
	panic(err)
}

for {
	row, err := reader.Read()

	if err == io.EOF {
		break
	}

	if err != nil {
		return err
	}

	value, ok := row["some-key"]

	// and so on...
}
```

### Writing files

```
import (
	csv "github.com/whosonfirst/go-whosonfirst-csv"
	"os"
)

fieldnames := []string{"foo", "bar"}

writer, err := csv.NewDictWriter(os.Stdout, fieldnames)

// or maybe you might do
// writer, err := csv.NewDictWriterFromPath("new.csv", fieldnames)

if err != nil {
	panic(err)
}

writer.WriteHeader()

row := make(map[string]string)
row["foo"] = "hello"
row["bar"] = "world"

// See this? "baz" is not included in the list of fieldnames
// above so it will be silently ignored and excluded from your
// CSV file. Perhaps it should trigger an error. It doesn't, today...

row["baz"] = "wub wub wub"

writer.WriteRow(row)
```

## Tools

### wof-csv-filter

Concatenate one or more CSV files ensuring a common set of columns.

```
./bin/wof-csv-filter -h
Usage of wof-csv-filter:
  $> wof-csv-filter -options <files>

Valid options are:
  -columns string
    	   Columns to filter on. A value of "-" means the set of unique columns for all CSV files being filtered. (default "-")
  -out string
       Where to write the data. A value of "-" means write the data to STDOUT. (default "-")
```

## See also

* https://golang.org/pkg/encoding/csv/
* https://docs.python.org/2/library/csv.html