package main

import (
	"fmt"

	"github.com/codercodingthecode/daix/currencyFetch/socketsclient"

	"github.com/codercodingthecode/daix/currencyFetch/currencyfetch"
)

func main() {
	ch := make(chan interface{})

	cL := currencyfetch.FetchCurrencies()

	// fmt.Println(cL)

	for k, v := range cL {
		fmt.Println(k)
		go socketsclient.SocketStream(k, v, ch)
	}

	// go socketsclient.SocketStream(cL, ch)
	for {
		fmt.Println("Websocket CH Return -> ", <-ch)
		<-ch
	}
}
