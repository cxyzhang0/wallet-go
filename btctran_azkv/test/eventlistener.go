// ref: https://yalantis.com/blog/how-to-build-websockets-in-go/
// ref: https://golangdocs.com/golang-gorilla-websockets
// ref: https://tradermade.com/tutorials/golang-websocket-client/
package test

import (
	"encoding/json"
	"fmt"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"testing"
	"time"
)

var (
	done      chan interface{}
	interrupt chan os.Signal
)

func receiveHandler(connection *websocket.Conn, t *testing.T) {
	defer close(done)
	count := 0
	for {
		msgType, msg, err := connection.ReadMessage()
		if err != nil {
			t.Errorf("error in receiveHandler.ReadMessage: %+v", err)
			return
		}
		t.Logf("receivied msg type %d: %s\n", msgType, msg)
		count++
		if count > 5 {
			connection.Close()
			return
		}
	}
}

func startWebsocketClient(t *testing.T, fromAddr string, txHash string) error {
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	//path := fmt.Sprintf("%s/%s/%s?token=%s", tran.Conf.Blockcypher.WebsocketPath, tran.Conf.Blockcypher.Coin, tran.Conf.Blockcypher.Chain, tran.Conf.Blockcypher.Token)
	//u := url.URL{Scheme: "wss", Host: tran.Conf.Blockcypher.WebsocketHost, Path: path}
	connStr := "wss://socket.blockcypher.com/v1/btc/test3?token=e905d13ae51748e2b618da1ba4ce0458" //u.String()
	//connStr := u.String() //
	t.Logf("connecting to %s", connStr)
	conn, httpResp, err := websocket.DefaultDialer.Dial(connStr, nil)
	if err != nil { // the above failed due to blockcypher.com certificate expired
		t.Fatalf("fail to dial websocket: %+v", err)
	}
	t.Logf("got dialer http resp: %+v", httpResp)

	defer conn.Close()

	// send event request
	event := kmssdk.TxConfirmationEvent{
		Event:         "tx-confirmation",
		Address:       fromAddr,
		Hash:          txHash,
		Confirmations: 6,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event to json: %+v", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, eventBytes)
	if err != nil {
		return fmt.Errorf("failed to write message to websocket: %+v", err)
	}

	go receiveHandler(conn, t)

	for {
		select {
		case <-done:
			t.Logf("all messages have been received. close web socket.")
			return nil
		case <-interrupt:
			t.Logf("interrupt")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				t.Logf("error write clode: %+v", err)
				return err
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
