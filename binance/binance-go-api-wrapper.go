/*
Copyright 2017 Wayland Chan
Licensed under terms of MIT license (see LICENSE)
*/

/*
Wrapper for the Binance Exchange REST API
*/

package binance

import (
	"net/http"
	"fmt"
	"strings"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"strconv"
	"errors"
)


const baseURL = "https://www.binance.com/api"

func init() {
	// do nothing
}

func New(myKey string, mySecret string) *BSession {
	return &BSession{httpClient: http.DefaultClient, key: myKey, secret: mySecret}
}

func (session *BSession) do(method, resource, payload string, auth bool, result interface{}) (resp *http.Response, err error) {

	fullUrl := fmt.Sprintf("%s/%s", baseURL, resource)

	req, err := http.NewRequest(method, fullUrl, strings.NewReader(payload))
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	if auth {
		if len(session.key) == 0 || len(session.secret) == 0 {
			err = errors.New("private endpoints require you to set an API Key and API Secret")
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
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		//uncomment next 2 lines if you need to see message for debugging
		//bodyString := string(bodyBytes)
		//fmt.Println("DEBUG:: " + bodyString)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&result)
	}
	return
}

func handleError(resp *http.Response) error {
	//Assuming anything other than HTTP 200 is an error. API docs don't really say
	if resp.StatusCode == 200 {
		return nil
	} else {
	body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		errorText := fmt.Sprintf("ERROR: %s", string(body))
		return errors.New(errorText)
	}

}

func (session *BSession) GetAccount(recv int) (acct Account, err error) {
	timestamp := time.Now().Unix() * 1000
	reqUrl := fmt.Sprintf("v3/account?timestamp=%d", timestamp)
	if recv > 0 && recv < 501 {
		reqUrl = fmt.Sprintf(reqUrl+"&recvWindow=%d", recv)
	}
	a := new(Account)
	_, err = session.do("GET", reqUrl, "", true, &a)

	return *a, err
}

func (session *BSession) GetOpenOrders(symbol string, receiveWindow int) (orders []Order, err error) {
	if symbol == "" {
		return nil, errors.New("symbol parameter is missing")
	}

	timestamp := time.Now().Unix() * 1000

	reqUrl := fmt.Sprintf("v3/openOrders?symbol=%s&timestamp=%d", symbol, timestamp)
	if receiveWindow > 0 {
		reqUrl = fmt.Sprintf(reqUrl+"&recvWindow=%d", receiveWindow)
	}

	orders = []Order{}
	_, err = session.do("GET", reqUrl, "", true, &orders)

	return orders, err
}

func (session *BSession) Get24Hr(symbol string) (price PriceChangeResponse, err error) {
	reqUrl := fmt.Sprintf("v1/ticker/24hr?symbol=%s", symbol)
	result := new(PriceChangeResponse)
	_, err = session.do("GET", reqUrl, "", false, &result)
	result.Symbol = symbol
	return *result, err
}

func (session *BSession) GetAllPrices() (prices []PriceTicker, err error) {
	reqUrl := "v1/ticker/allPrices"
	results := []PriceTicker{}
	_, err = session.do("GET", reqUrl, "", false, &results)

	return results, err
}

func (session *BSession) GetOrderBook(symbol string, limit int) (ob OrderBook, err error) {
	reqUrl := fmt.Sprintf("v1/depth?symbol=%s&limit=%d", symbol, limit)
	result := new(OrderBook)
	_, err = session.do("GET", reqUrl, "", false, &result)
	result.Symbol = symbol

	return *result, err
}

func (session *BSession) CreateOrder(req OrderRequest) (res OrderResponse, err error) {
	var reqUrl string
	if req.Type == "LIMIT" {
		reqUrl = fmt.Sprintf("v3/order?symbol=%s&type=%s&side=%s&timeInForce=%s&quantity=%f&price=%f",
			req.Symbol, req.Type, req.Side, req.TimeInForce, req.Quantity, req.Price)
	} else if req.Type == "MARKET" {
		reqUrl = fmt.Sprintf("v3/order?symbol=%s&type=%s&side=%s&quantity=%f",
			req.Symbol, req.Type, req.Side, req.Quantity)

	}

	result := new(OrderResponse)
	_, err = session.do("POST", reqUrl, "", true, &result)
	result.Symbol = req.Symbol
	return *result, err
}

func (session *BSession) QueryOrder(req OrderResponse) (status OrderStatus, err error) {
	reqUrl := fmt.Sprintf("v3/order?symbol=%s&orderId=%d", req.Symbol, req.OrderId)
	result := new(OrderStatus)
	_, err = session.do("GET", reqUrl, "", true, &result)
	return *result, err
}

func (session *BSession) CancelOrder(req OrderResponse) (res CancelOrderResponse, err error) {
	reqUrl := fmt.Sprintf("v3/order?symbol=%s&orderId=%d&origClientOrderId=%s",
		req.Symbol, req.OrderId, req.ClientOrderId)
	result := new(CancelOrderResponse)
	_, err = session.do("DELETE", reqUrl, "", true, &result)
	return *result, err
}

func (o *OrderQuote) UnmarshalJSON(b []byte) error {
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
