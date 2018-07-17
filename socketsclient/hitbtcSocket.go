package socketsclient

import (
	"encoding/json"
	"fmt"

	ws "github.com/gorilla/websocket"
)

// const websckt = "wss://api.hitbtc.com/api/2/ws"
const gdax = "wss://ws-feed.gdax.com"
const method = "getSymbol"
const params = "ETHBTC"

// WSHitbtc is pulling data
func WSHitbtc(cL map[string][]string, ch chan interface{}) {

	// fmt.Println("DATA COMING FROM MAIN ---------------> ", cL["gdax"])

	var sliceData []string
	for _, v := range cL["gdax"] {
		if v != "" {
			sliceData = append(sliceData, v)
		}
	}

	fmt.Println("after shaped -> ", sliceData)

	var websocket ws.Dialer

	wsConn, _, wsError := websocket.Dial(gdax, nil)

	if wsError != nil {
		fmt.Println("There was an error dialing in -> ", wsError)
	}

	fmt.Println("------------------------------------------------------------")
	// fmt.Println("This is wsConn -> ", wsConn.)
	b, _ := json.Marshal(sliceData)
	bigString := fmt.Sprintf(`{"type":"subscribe","channels": [{"name": "ticker", "product_ids": %s}]}`, b)

	var writeData interface{}

	writeDataError := json.Unmarshal([]byte(bigString), &writeData)
	if writeDataError != nil {
		fmt.Println("There was an eror unmarshalling ", writeDataError)
	}

	fmt.Println("writedata json ", writeData)

	writeJSONError := wsConn.WriteJSON(writeData)
	if writeJSONError != nil {
		fmt.Println("there was an error writing json -> ", writeJSONError)
	}

	var tickRawData interface{}
	for {
		// break
		readError := wsConn.ReadJSON(&tickRawData)

		if readError != nil {
			fmt.Println("Error in pulling in data -> ", readError)
		}

		// tickData := tickRawData.(map[string]interface{})

		// fmt.Println("This is the feed -> ", tickData)
		ch <- tickRawData
	}

}
