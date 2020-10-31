package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type doneCode struct {
	ExitCode int
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:28114"}
	log.Printf("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial error: ", err)
	}
	defer c.Close()

	done := make(chan doneCode)

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"id": 2, "jsonrpc": "2.0", "method": "subscribe", "params": ["new_tip_header"]}`))
	if err != nil {
		log.Println("write error: ", err)
	}

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				done <- doneCode{1}
				return
			}
			log.Printf("recv: %s\n", message)
		}
	}()

	for {
		select {
		case doneCode := <-done:
			log.Println("done.")
			os.Exit(doneCode.ExitCode)
		case <-interrupt:
			log.Println("interrupted.")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye."))
			if err != nil {
				log.Println("write close error: ", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}