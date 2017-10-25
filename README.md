# Go wrapper for Binance REST API

binance-go-api-wrapper is a Go client library wrapper for accessing Binance's web REST API.

For details of the REST API, see Binance's website
https://www.binance.com/restapipub.html

Initially created as a submission to their API coding competition but I'll continue working on this for my own eventual use.
https://support.binance.com/hc/en-us/articles/115001909972

## Install ##
Install the package with:
```bash
go get github.com/waylandc/binance-go-api-wrapper/binance
```

## Example ##
Before running this, you need to create API keys from Binance website
and set the BINANCE_KEY and BINANCE_SECRET environment variables.

```go
package main
import (
	"os"
	"fmt"
	"github.com/waylandc/binance-go-api-wrapper/binance"
)

func main() {
	session := binance.New(os.Getenv("BINANCE_KEY"), os.Getenv("BINANCE_SECRET"))
	prices, err := session.GetAllPrices()
	if err != nil {
		fmt.Println(err)
	}

	i := 0
	for i < len(prices)-1 {
		fmt.Printf("Symbol: %s Price: %f\n", prices[i].Symbol, prices[i].Price)
		i++
	}
	
}
```

## License ##
MIT License - see LICENSE file for details
