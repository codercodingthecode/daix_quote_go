// THIS SNIPPET FETCHS THE EXISTING CURRENCIES FROM THE EXCHANGES THRU REST API CALLS
// THE PURPOSE IS TO BYPASS THE CALL TO DYNAMODB FROM WEBSOCKETS CODE AND HAVE A MORE ACCURATE LIST OF CURRENCY PER EXCHANGE.

//@todo:
// handle errors better
// export cL variable holding the current data from exchanges

package currencyfetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// CoinRest interface holder for rest api incoming data
type CoinRest interface{}

// FetchCurrencies return a list with all pairs per exchange
func FetchCurrencies() map[string][]string {
	// general use channel to receive async from functions calling the rest api's
	var c = make(chan []string)

	// cL -> "coin list" contains a dictionary of array containing the pairs within it.
	// It holds the exchange as key and the currencies supported on that specific exchange as a slice
	// ie:
	// exchange: currency
	// gdax: ["BTC", "EUR", "LTC", "GBP", "USD", "ETH", "BCH"]
	// data structure that holds the whole exchanges supported currecy list.
	var cL = make(map[string][]string)

	// list of the current exchange rest api's to fetch the currencies from each exchange.
	// it currently calls the exchanges that provides websockets, so it can provide a list of currencies
	// to the websockets api per exchange.
	// urls := map[string]string{
	// 	// "hitbtc": "https://api.hitbtc.com/api/2/public/symbol", // response -> id -> "ETHUSD" **
	// 	"gdax": "https://api.gdax.com/products", // response -> id -> "eth" **
	// 	// "binance":  "https://api.binance.com/api/v1/exchangeInfo",          // it needs to dig into json, response -> symbol -> baseAsset -> "ETH" **
	// 	// "poloniex": "https://poloniex.com/public?command=returnCurrencies", // response -> "ETH"   ; the crypto is in the key **
	// 	// https://poloniex.com/public?command=returnTicker  *change to this end point
	// 	// "bitmex":   "https://www.bitmex.com/api/v1/stats", // needs to look into it further to find pairs                 // response -> rootSymbol -> "ETH" **
	// 	// "bitfinex": "https://api.bitfinex.com/v1/symbols", // response -> "ethbtc"  || it takes a pair rather than the currency symbol, needs discussion on it
	// 	// "bittrex": "https://bittrex.com/api/v1.1/public/getcurrencies", // response -> result -> Currency -> "ETH" || it uses rest api already.
	// }

	// range over the rest api list and call the api's function.
	for k, v := range exchangeList() {
		switch k {
		case "hitbtc":
			go fethHitbtc(v, c)
			cL[k] = <-c
		case "gdax":
			go fethGdax(v, c)
			cL[k] = <-c
		case "bitmex":
			go fetchBitmex(v, c)
			cL[k] = <-c
		case "poloniex":
			go fetchPoloniex(v, c)
			cL[k] = <-c
		case "binance":
			go fetchBinance(v, c)
			cL[k] = <-c
		}

	}

	// holds the final data with exchanges and currencies
	// @todo: needs to change it to be exported and consumed by other code.
	// fmt.Println("Currency List ->> ", cL)
	// socketsclient.WSHitbtc(cL)

	// testing new code. return the entire list of pairs per exchange
	return cL
}

// Helper function common used to by the function below to call the rest api for each exchange.
// Function: Receives two arguments and returns the data by pointer to the caller:
// - url string which is the exchange rest api
// - coins pointer interface which receives the data from the rest call and returns to the caller function.
func restAPICaller(url string, coins *CoinRest) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("there was an error -> ", err)
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(bytes, &coins)
}

// Fetchs the currency list from POLONIEX thru rest call
// Function: Receives two arguemts and feed data thru the channel back to the caller:
// - url string which is the exchange rest api
// - c channel which is where it feeds the data back to the caller once ready
// it call "restAPICaller" function passing the rest url and a coin interface pointer
// to retrieve the data
func fetchPoloniex(url string, c chan []string) {
	cList := make([]string, 100)
	var coins CoinRest

	restAPICaller(url, &coins)

	for i := range coins.(map[string]interface{}) {
		cList = append(cList, i)
	}
	c <- cList
}

// Fetchs the currency list from BITMEX thru rest call
// Function: Receives two arguemts and feed data thru the channel back to the caller:
// - url string which is the exchange rest api
// - c channel which is where it feeds the data back to the caller once ready
// it call "restAPICaller" function passing the rest url and a coin interface pointer
// to retrieve the data
func fetchBitmex(url string, c chan []string) {
	cList := make([]string, 100)
	var coins CoinRest

	restAPICaller(url, &coins)

	for _, v := range coins.([]interface{}) {
		s := v.(map[string]interface{})
		cList = append(cList, s["rootSymbol"].(string))
	}
	c <- cList
}

// Fetchs the currency list from BINANCE thru rest call
// Function: Receives two arguemts and feed data thru the channel back to the caller:
// - url string which is the exchange rest api
// - c channel which is where it feeds the data back to the caller once ready
// it call "restAPICaller" function passing the rest url and a coin interface pointer
// to retrieve the data
func fetchBinance(url string, c chan []string) {
	cList := make([]string, 100)
	var coins CoinRest

	restAPICaller(url, &coins)

	for i, v := range coins.(map[string]interface{}) {
		if i == "symbols" {
			for _, coin := range v.([]interface{}) {
				s := coin.(map[string]interface{})
				cList = append(cList, s["baseAsset"].(string))
			}
		}
	}
	c <- cList
}

// Fetchs the currency list from HITBTC thru rest call which had the same data structure
// Function: Receives two arguemts and feed data thru the channel back to the caller:
// - url string which is the exchange rest api
// - c channel which is where it feeds the data back to the caller once ready
// it call "restAPICaller" function passing the rest url and a coin interface pointer
// to retrieve the data
func fethHitbtc(url string, c chan []string) {
	cList := make([]string, 900)
	var coins CoinRest

	restAPICaller(url, &coins)

	for _, coin := range coins.([]interface{}) {
		s := coin.(map[string]interface{})
		cList = append(cList, s["id"].(string))
	}
	c <- cList
}

// Fetchs the currency list from GDAX thru rest call which had the same data structure
// Function: Receives two arguemts and feed data thru the channel back to the caller:
// - url string which is the exchange rest api
// - c channel which is where it feeds the data back to the caller once ready
// it call "restAPICaller" function passing the rest url and a coin interface pointer
// to retrieve the data
func fethGdax(url string, c chan []string) {
	cList := make([]string, 900)
	var coins CoinRest

	restAPICaller(url, &coins)

	for _, coin := range coins.([]interface{}) {
		s := coin.(map[string]interface{})
		cList = append(cList, s["id"].(string))
	}
	c <- cList
}
