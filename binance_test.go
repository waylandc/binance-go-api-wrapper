package main
import (
	"os"
	"testing"
	"net/http"
)

var sess *MySession

func TestMain(m *testing.M) {
	//sess = binance.New(os.Getenv("BINANCE_KEY"), os.Getenv("BINANCE_SECRET"))
	sess = &MySession{
		httpClient: new(http.Client),
		key: os.Getenv("BINANCE_KEY"),
		secret: os.Getenv("BINANCE_SECRET"),
	}

	code := m.Run()
	os.Exit(code)
}

func Test24HourPrice(t *testing.T) {
	price, err := sess.get24Hr("BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Symbol: %v Bid: %v Ask: %v Last: %v Volume: %v\n",
		price.Symbol, price.BidPrice, price.AskPrice, price.LastPrice, price.Volume)

}