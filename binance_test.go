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
	price, err := sess.Get24Hr("BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Symbol: %v Bid: %v Ask: %v Last: %v Volume: %v\n",
		price.Symbol, price.BidPrice, price.AskPrice, price.LastPrice, price.Volume)

}

func TestGetAllPrices(t *testing.T) {
	prices := []PriceTicker{}
	prices, err := sess.GetAllPrices()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Returned %v prices", len(prices))
	i := 0
	for i < 3 {
		t.Logf("Symbol: %v Price: %v\n", prices[i].Symbol, prices[i].Price)
		i++
	}
}

func TestGetOrderBook(t *testing.T) {
	ob, err := sess.GetOrderBook("BTCUSDT", 10)
	if err != nil {
		t.Error(err)
	}

	if ob.Symbol != "BTCUSDT" {
		t.Errorf("Expected order book for %s. Got %s", "BTCUSDT", ob.Symbol)
	}

	t.Logf("Order book for %s returned %d asks and %d bids", ob.Symbol, len(ob.Asks), len(ob.Bids))
}