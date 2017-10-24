package main

/**
	Wrapper for the Binance Exchange REST API
 */

import (
	"net/http"
	"strings"
	"fmt"
	"time"
	"errors"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

const baseURL = "https://www.binance.com/api"

type MySession struct {
	httpClient *http.Client
	key string
	secret string
}

type PriceChangeResponse struct {
	Symbol string
	PriceChange string	`json:"priceChange"`
	PriceChangePercent string	`json:"priceChangePercent"`
	WeightedAvgPrice float64	`json:"weightedAvgPrice,string"`
	PrevClosePrice float64	`json:"prevClosePrice,string"`
	LastPrice float64	`json:"lastPrice,string"`
	BidPrice float64	`json:"bidPrice,string"`
	AskPrice float64	`json:"askPrice,string"`
	OpenPrice float64	`json:"openPrice,string"`
	HighPrice float64	`json:"highPrice,string"`
	LowPrice float64	`json:"lowPrice,string"`
	Volume float64		`json:"volume,string"`
	OpenTime int64		`json:"openTime"`
	CloseTime int64		`json:"closeTime"`
	Count int64			`json:"count"`
	// don't know what these are used for so comment out but leave here as reminder they're available
	//FirstId int64	`json:"firstId"`
	//LastId int64	`json:"lastId"`
}

type PriceTicker struct {
	Symbol string	`json:"symbol"`
	Price float64	`json:"price,string"`
}

type Order struct {
	Price float64	`json:",string"`
	Quantity float64	`json:",string"`
}

type OrderBook struct {
	Symbol string
	LastUpdateId int	`json:"lastUpdateId"`
	Bids []Order 	`json:",string"`
	Asks []Order	`json:",string"`
}

func New(myKey string, mySecret string) *MySession {
	return &MySession{httpClient:http.DefaultClient, key: myKey, secret: mySecret}
}

func (session *MySession) do(method, resource, payload string, auth bool, result interface{}) (resp *http.Response, err error) {

	fullUrl := fmt.Sprintf("%s/%s", baseURL, resource)

	req, err := http.NewRequest(method, fullUrl, strings.NewReader(payload))
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	if auth {
		if len(session.key) == 0 || len(session.secret) == 0 {
			err = errors.New("Private endpoints requre you to set an API Key and API Secret")
			return
		}

		req.Header.Add("X-MBX-APIKEY", session.key)

		q := req.URL.Query()

		timestamp := time.Now().Unix() * 1000
		q.Set("timestamp", fmt.Sprintf("%d", timestamp))

		mac := hmac.New(sha256.New, []byte(session.secret))
		_, err := mac.Write([]byte(q.Encode()))
		if err != nil {
			return nil, err
		}

		signature := hex.EncodeToString(mac.Sum(nil))
		req.URL.RawQuery = q.Encode() + "&signature=" + signature
	}

	resp, err = session.httpClient.Do(req)
	if err != nil {
		return
	}

	// Check for error
	err = handleError(resp)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	// Process response
	if resp != nil {
		//bodyBytes, _ := ioutil.ReadAll(resp.Body)
		//bodyString := string(bodyBytes)
		//fmt.Println(bodyString)
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&result)
	}
	return
}

func handleError(resp *http.Response) error {

	if resp.StatusCode == 400 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		errorText := fmt.Sprintf("Bad Request: %s", string(body))
		return errors.New(errorText)
	} else {
		return nil
	}

}

//func (session *Session) getOpenOrders() {
//
//}
//
//func getOrderBook() {
//
//}

func (session *MySession) Get24Hr(symbol string) (price PriceChangeResponse, err error) {
	reqUrl := fmt.Sprintf("v1/ticker/24hr?symbol=%s", symbol)
	result := new(PriceChangeResponse)
	_, err = session.do("GET", reqUrl, "", false, &result)
	result.Symbol = symbol
	return *result, err
}

func (session *MySession) GetAllPrices() (prices []PriceTicker, err error) {
	reqUrl := "v1/ticker/allPrices"
	results := []PriceTicker{}
	_, err = session.do("GET", reqUrl, "", false, &results)

	return results, err
}

func (session *MySession) GetOrderBook(symbol string, limit int) (ob OrderBook, err error) {
	reqUrl := fmt.Sprintf("v1/depth?symbol=%s&limit=%d", symbol, limit)
	result := new(OrderBook)
	_, err = session.do("GET", reqUrl, "", false, &result)
	result.Symbol = symbol

	return *result, err
}

func (o *Order) UnmarshalJSON(b []byte) error {
	var s [2]string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	o.Price, err = strconv.ParseFloat(s[0], 64)
	if err != nil {
		return err
	}

	o.Quantity, err = strconv.ParseFloat(s[1], 64)
	if err != nil {
		return err
	}

	return nil
}