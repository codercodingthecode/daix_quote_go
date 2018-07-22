package main

import (
	"fmt"

	"github.com/codercodingthecode/daix/currencyFetch/socketsclient"

	"github.com/codercodingthecode/daix/currencyFetch/currencyfetch"
)

func main() {
	ch := make(chan interface{})
	// chOut := make(chan interface{})
	cL := currencyfetch.FetchCurrencies()

	for k, v := range cL {
		fmt.Println("RUNNING ON MAIN")
		go socketsclient.SocketStream(k, v, ch)
	}

	for {
		fmt.Println(<-ch)
	}
	// exchanges.Gdax(ch)

}
