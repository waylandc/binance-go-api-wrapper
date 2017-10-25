/*
Copyright 2017 Wayland Chan
Licensed under terms of MIT license (see LICENSE)
*/

package binance

import (
	"net/http"
	"os"
	"testing"
)

var session *BSession

const TestSymbol = "BTCUSDT"

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestSetup(t *testing.T) {
	if os.Getenv("BINANCE_KEY") == "" || os.Getenv("BINANCE_SECRET") == "" {
		print("BINANCE_KEY AND BINANCE_SECRET env vars must be set")
		os.Exit(1)
	}

	session = &BSession{
		httpClient: new(http.Client),
		key:        os.Getenv("BINANCE_KEY"),
		secret:     os.Getenv("BINANCE_SECRET"),
	}
}

func Test24HourPrice(t *testing.T) {
	t.Log("Get the latest price of " + TestSymbol)
	price, err := session.Get24Hr(TestSymbol)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Symbol: %v Bid: %v Ask: %v Last: %v Volume: %v\n",
		price.Symbol, price.BidPrice, price.AskPrice, price.LastPrice, price.Volume)

}

func TestGetOrderBook(t *testing.T) {
	t.Log("Get the depth/order book of " + TestSymbol)
	ob, err := session.GetOrderBook(TestSymbol, 10)
	if err != nil {
		t.Error(err)
	}

	if ob.Symbol != TestSymbol {
		t.Errorf("Expected order book for %s. Got %s", TestSymbol, ob.Symbol)
	}

	t.Logf("OrderQuote book for %s returned %d asks and %d bids", ob.Symbol, len(ob.Asks), len(ob.Bids))
}

func TestCreateQueryCancelLimitOrder(t *testing.T) {
	t.Log("Place a LIMIT order")
	// Create order request and submit
	req := &OrderRequest{
		Symbol:      TestSymbol,
		Side:        "SELL",
		Type:        "LIMIT",
		TimeInForce: "GTC",
		Quantity:    0.001,
		Price:       9500.00}

	res, err := session.CreateOrder(*req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Limit Ordered ID %d submitted, %s", res.OrderId, res.ClientOrderId)

	t.Log("Check an order's status")
	// Query the order we just submitted
	query, err := session.QueryOrder(res)
	if err != nil {
		t.Error(err)
	}

	t.Logf("OrderQuery::\nOrder ID=%d, Symbol=%s, Side=%s, Status=%s, Price=%f, Quantity=%f, ExecutedQty=%f",
		query.OrderId, query.Symbol, query.Side, query.Status, query.Price, query.OrigQuantity, query.ExecutedQuantity)

	t.Log("Cancel an order")
	//Cancel the same order
	cancel, err := session.CancelOrder(res)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Cancelled order ID %d, cancel reference %s", cancel.OrderId, cancel.ClientOrderId)
}

func TestCreateMarketOrder(t *testing.T) {
	t.Log("Place a MARKET order")
	req := &OrderRequest{
		Symbol:"STRATBTC",
		Side:"BUY",
		Type:"MARKET",
		TimeInForce:"GTC",
		Quantity:10,
	}

	res, err := session.CreateOrder(*req)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Market Ordered ID %d submitted, %s", res.OrderId, res.ClientOrderId)

	t.Log("Check an order's status")
	// Query the order we just submitted
	query, err := session.QueryOrder(res)
	if err != nil {
		t.Error(err)
	}

	t.Logf("OrderQuery::\nOrder ID=%d, Symbol=%s, Side=%s, Status=%s, Price=%f, Quantity=%f, ExecutedQty=%f",
		query.OrderId, query.Symbol, query.Side, query.Status, query.Price, query.OrigQuantity, query.ExecutedQuantity)

}

//func TestGetAllPrices(t *testing.T) {
//	t.Log("Get all prices")
//	prices, err := session.GetAllPrices()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	t.Logf("Returned %v prices", len(prices))
//	i := 0
//	for i < 3 {
//		t.Logf("Symbol: %v Price: %v\n", prices[i].Symbol, prices[i].Price)
//		i++
//	}
//}


func TestGetOpenOrders(t *testing.T) {
	t.Log("Get list of open orders")
	openOrders, err := session.GetOpenOrders(TestSymbol, 0)
	if err != nil {
		t.Error(err)
	}
	t.Logf("We currently have %d open orders.", len(openOrders))

}

func TestGetAccount(t *testing.T) {
	t.Log("Get account info and show list of all current positions")
	acct, err := session.GetAccount(0)
	if err != nil {
		t.Error(err)
	}

	t.Log("Account Summary::")
	t.Logf("\tMaker Commission: %2d%%\tTaker Commission:  %2d%% ", acct.MakerCommission, acct.TakerCommission)
	t.Logf("\tBuyer Commission: %2d%%\tSeller Commission: %2d%% ", acct.BuyerCommission, acct.SellerCommission)
	t.Logf("\tCan Trade: %v", acct.CanTrade)
	t.Logf("\tCan Deposit: %v", acct.CanDeposit)
	t.Logf("\tCan Withdraw: %v", acct.CanWithdraw)
	t.Log("\tPositions::")
	t.Log("\t    Coin\t Available \t In Order")
	for _, bal := range acct.Balances {
		t.Logf("\t\t%s  \t %f  \t %f", bal.Asset, bal.Free, bal.Locked)
	}
}
