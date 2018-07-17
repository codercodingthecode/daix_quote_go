package socketsclient

import (
	"encoding/json"
	"fmt"

	ws "github.com/gorilla/websocket"
)

// SocketStream common function for all exchange, returns raw data
func SocketStream(cL map[string][]string, ch chan interface{}) {
	var writeData interface{}
	var tickRawData interface{}
	var websocket ws.Dialer
	wsConn, _, wsError := websocket.Dial(WebsocketList()["gdax"], nil) // this will have to change per loop or concurrency
	if wsError != nil {
		fmt.Println("There was an error dialing in -> ", wsError)
	}

	b, _ := json.Marshal(cL["gdax"])

	// move this to struct
	bigString := fmt.Sprintf(`{"type":"subscribe","channels": [{"name": "ticker", "product_ids": %s}]}`, b)

	writeDataError := json.Unmarshal([]byte(bigString), &writeData)
	if writeDataError != nil {
		fmt.Println("There was an eror unmarshalling ", writeDataError)
	}
	writeJSONError := wsConn.WriteJSON(writeData)
	if writeJSONError != nil {
		fmt.Println("there was an error writing json -> ", writeJSONError)
	}
	for {
		readError := wsConn.ReadJSON(&tickRawData)
		if readError != nil {
			fmt.Println("Error in pulling in data -> ", readError)
		}

		// tickData := tickRawData.(map[string]interface{})

		// fmt.Println("This is the feed -> ", tickData)
		ch <- tickRawData
	}

}
