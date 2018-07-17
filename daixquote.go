package main

import (
	"fmt"

	"github.com/codercodingthecode/daix/currencyFetch/currencyfetch"
	"github.com/codercodingthecode/daix/currencyFetch/socketsclient"
)

func main() {
	ch := make(chan interface{})
	cL := currencyfetch.FetchCurrencies()
	go socketsclient.WSHitbtc(cL, ch)
	fmt.Println("Websocket CH Return -> ", <-ch)
}
