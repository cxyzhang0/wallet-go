// ref: https://golangdocs.com/golang-gorilla-websockets
package test

import (
	"github.com/gorilla/websocket"
	"net/url"
	"os"
	"os/signal"
	"testing"
)

var (
	done      chan interface{}
	interrupt chan os.Signal
)

func receiveHandler(connection *websocket.Conn, t *testing.T) {
	defer close(done)
	for {
		msgType, msg, err := connection.ReadMessage()
		if err != nil {
			t.Errorf("error in receiveHandler.ReadMessage: %+v", err)
			return
		}
		t.Logf("receivied msg type %d: %s\n", msgType, msg)
	}
}

func startWebsocketClient(t *testing.T, addr string) {
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	u := url.URL{Scheme: "wss", Host: wsHost, Path: wsPath}
	t.Logf("connecting to %s", u.String())
	conn, httpResp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("fail to dial websocket: %+v", err)
	}
	t.Logf("got dialer http resp: %+v", httpResp)

	defer conn.Close()

	// send event request
	conn.WriteMessage(websocket.TextMessage, []byte(""))

	go receiveHandler(conn, t)

}
