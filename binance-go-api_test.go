package main
import (
	"os"
	"testing"
	"net/http"
)

var sess *MySession
const TEST_SYMBOL = "BTCUSDT"

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
	price, err := sess.Get24Hr(TEST_SYMBOL)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Symbol: %v Bid: %v Ask: %v Last: %v Volume: %v\n",
		price.Symbol, price.BidPrice, price.AskPrice, price.LastPrice, price.Volume)

}

func TestGetAllPrices(t *testing.T) {
	//prices := []PriceTicker{}
	prices, err := sess.GetAllPrices()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Returned %v prices", len(prices))
	i := 0
	for i < 3 {
		t.Logf("Symbol: %v Price: %v\n", prices[i].Symbol, prices[i].Price)
		i++
	}
}

func TestGetOrderBook(t *testing.T) {
	ob, err := sess.GetOrderBook(TEST_SYMBOL, 10)
	if err != nil {
		t.Error(err)
	}

	if ob.Symbol != "BTCUSDT" {
		t.Errorf("Expected order book for %s. Got %s", TEST_SYMBOL, ob.Symbol)
	}

	t.Logf("OrderQuote book for %s returned %d asks and %d bids", ob.Symbol, len(ob.Asks), len(ob.Bids))
}

func TestGetAccount(t *testing.T) {
	acct, err := sess.GetAccount(0)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Your account can trade, %v", acct.CanTrade)
}

func TestCreateLimitOrder(t *testing.T) {
	req := &OrderRequest{
		Symbol:TEST_SYMBOL,
		Side:"SELL",
		Type:"LIMIT",
		TimeInForce:"GTC",
		Quantity:0.001,
		Price:9500.00,}

	res, err := sess.CreateOrder(*req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Ordered ID %d submitted, %s", res.OrderId, res.ClientOrderId)

}

func TestGetOpenOrders(t *testing.T) {
	openOrders, err := sess.GetOpenOrders(TEST_SYMBOL, 0)
	if err != nil {
		t.Error(err)
	}
	t.Logf("We currently have %d open orders.", len(openOrders))

}

