package binance

import "net/http"

type BSession struct {
	httpClient *http.Client
	key        string
	secret     string
}

type Account struct {
	MakerCommission  int       `json:"makerCommission"`
	TakerCommission  int       `json:"takerCommission"`
	BuyerCommission  int       `json:"buyerCommission"`
	SellerCommission int       `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	Balances         []Balance `json:",string"`
}

type Balance struct {
	Asset  string  `json:"asset"`         //symbol of coin/asset
	Free   float64 `json:"free,string"`   //available balance
	Locked float64 `json:"locked,string"` //quantity in trades
}

type PriceChangeResponse struct {
	Symbol             string
	PriceChange        string  `json:"priceChange"`
	PriceChangePercent string  `json:"priceChangePercent"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
	PrevClosePrice     float64 `json:"prevClosePrice,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	BidPrice           float64 `json:"bidPrice,string"`
	AskPrice           float64 `json:"askPrice,string"`
	OpenPrice          float64 `json:"openPrice,string"`
	HighPrice          float64 `json:"highPrice,string"`
	LowPrice           float64 `json:"lowPrice,string"`
	Volume             float64 `json:"volume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	Count              int64   `json:"count"`
	// don't know what these are used for so comment out but leave here as reminder they're available
	//FirstId int64	`json:"firstId"`
	//LastId int64	`json:"lastId"`
}

type PriceTicker struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

// used in allOrders
type Order struct {
	Symbol           string  `json:"symbol"`
	OrderId          int     `json:"orderId"`
	ClientOrderId    string  `json:"clientOrderId"`
	Price            float64 `json:"price,string"`
	OrigQuantity     float64 `json:"origQty,string"`
	ExecutedQuantity float64 `json:"executedQty,string"`
	Status           string  `json:"status"`
	TimeInForce      string  `json:"timeInForce"`
	Type             string  `json:"type"`
	Side             string  `json:"side"`
	StopPrice        float64 `json:"stopPrice,string"`
	IcebergQuantity  float64 `json:"icebergQty,string"`
	Time             int64   `json:"time"`
}

type OrderQuote struct {
	Price    float64 `json:",string"`
	Quantity float64 `json:",string"`
}

type OrderBook struct {
	Symbol       string
	LastUpdateId int          `json:"lastUpdateId"`
	Bids         []OrderQuote `json:",string"`
	Asks         []OrderQuote `json:",string"`
}

type OrderRequest struct {
	Symbol           string  `json:"symbol"`
	Side             string  `json:"side"`
	Type             string  `json:"type"`
	TimeInForce      string  `json:"timeInForce"`
	Quantity         float64 `json:"quantity,string"`
	Price            float64 `json:"price,string"`
	NewClientOrderId string  `json:"newClientOrderId"`  // not mandatory. uniq id, auto generated
	StopPrice        float64 `json:"stopPrice,string"`  // not mandatory, used with stop orders
	IcebergQuantity  float64 `json:"icebergQty,string"` // not mandatory, used with iceberg orders
}

type OrderResponse struct {
	Symbol        string `json:"symbol"`
	OrderId       int64  `json:"orderId"`
	ClientOrderId string `json:"clientOrderId"`
	TransactTime  int64  `json:"transactTime"`
}

type OrderStatus struct {
	Symbol           string  `json:"symbol"`
	OrderId          int64   `json:"orderId"`
	ClientOrderId    string  `json:"clientOrderId"`
	Price            float64 `json:"price,string"`
	OrigQuantity     float64 `json:"origQty,string"`
	ExecutedQuantity float64 `json:"executedQty,string"`
	Status           string  `json:"status"`
	TimeInForce      string  `json:"timeInForce"`
	Type             string  `json:"type"`
	Side             string  `json:"side"`
	StopPrice        float64 `json:"stopPrice,string"`
	IcebergQuantity  float64 `json:"icebergQty,string"`
	Time             int64   `json:"time"`
}

type CancelOrderResponse struct {
	Symbol            string `json:"symbol"`
	OrigClientOrderId string `json:"origClientOrderId"`
	OrderId           int64  `json:"orderId"`
	ClientOrderId     string `json:"clientOrderId"`
}
