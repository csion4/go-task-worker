package ws

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var OverFlag = []byte("Over!")

// websocket客户端
func WebSocketClient(ch *chan []byte, taskCode string, recordId int) {
	header := http.Header{}
	header.Set("auth", os.Getenv("Auth"))
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws/clusterResp?taskCode=%s&recordId=%d", os.Getenv("MNode"), taskCode, recordId), header)
	if err != nil {
		log.Println("Error connecting to Websocket Server:", err)
		return
	}
	defer conn.Close()
	// go receiveHandler(conn)

	for ;; {
		logs := <- *ch
		err = conn.WriteMessage(websocket.TextMessage, logs)
		if err != nil {
			log.Println("Error WriteMessage to Websocket Server:", err)
			return
		}
		if bytes.Equal(logs, OverFlag) {
			fmt.Println(taskCode + " -> task over!")
			return
		}
	}
}

