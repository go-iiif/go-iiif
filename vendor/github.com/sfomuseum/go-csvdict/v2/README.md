# go-csvdict

Go package to implement a "dict reader" style CSV parser (on top of the default `encoding/csv` package) to return rows a key-value dictionaries rather than lists.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-csvdict.svg)](https://pkg.go.dev/github.com/sfomuseum/go-csvdict)

## Example

### Reading files

```
import (
	"os"

        "github.com/whosonfirst/go-csvdict/v2"
)

r, _ := csvdict.NewReaderFromPath("example.csv")

// or maybe you might do
// r, _ := csvdict.NewReader(os.Stdin)

for {
	row, err := r.Read()

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

It is also possible to iterate through all the records using the `Iterate` method:

```
import (
	"os"

        "github.com/whosonfirst/go-csvdict/v2"
)

r, _ := csvdict.NewReaderFromPath("example.csv")

// or maybe you might do
// r, _ := csvdict.NewReader(os.Stdin)

for row, err := r.Iterate() {

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
	"os"

	"github.com/whosonfirst/go-csvdict/v2"
)

wr, _ := csvdict.NewWriter(os.Stdout)

// or maybe you might do
// wrr, _ := csvdict.NewWriterFromPath("new.csv")

row := make(map[string]{
	"foo": "hello",
	"bar": "world",
}

wr.WriteRow(row)
```

## See also

* https://golang.org/pkg/encoding/csv/
* https://docs.python.org/2/library/csv.html