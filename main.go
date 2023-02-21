package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// http 요청을 websocket 요청으로 승격한다.
		conn, err := upgrader.Upgrade(w, r, nil)

		// error가 발생하면 Handler를 종료한다.
		if err != nil {
			log.Fatal(err)
		}

		// Handler함수가 종료될 때 connection을 종료한다.
		// client 연결이 종료되면 handler가 종료된다.
		defer func() {
			conn.Close()
			fmt.Println("connection closed")
		}()

		// websocket 통신 client connection으로 부터 메시지를 받아들인다.
		for {
			messageType, message, err := conn.ReadMessage()

			if err != nil {
				// client와 연결이 끊어져도 에러로 잡힌다.
				fmt.Println(err)
				return
			}

			var stringMessage string = fmt.Sprintf("%s", message)
			fmt.Println(messageType, stringMessage)

			if stringMessage == "나는" {
				conn.WriteMessage(messageType, []byte("바보야"))
			}

		}
	})

	fmt.Println("Listen server Port :9080")
	if err := http.ListenAndServe(":9080", nil); err != nil {
		log.Fatal(err)
	}
}
