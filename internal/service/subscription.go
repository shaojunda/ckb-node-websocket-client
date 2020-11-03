package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/shaojunda/ckb-node-websocket-client/global"
	"github.com/shaojunda/ckb-node-websocket-client/internal/rpc"
)

var supportedTopics = []string{"new_tip_header", "new_tip_block", "new_transaction"}

func (svc Service) Subscribe(c *websocket.Conn, topic string) error {
	if !contains(supportedTopics, topic) {
		return fmt.Errorf("topic %s is not supported", topic)
	}
	err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"id": 2, "jsonrpc": "2.0", "method": "subscribe", "params": ["%s"]}`, topic)))
	if err != nil {
		global.Logger.Errorf("write error: ", err)
	}
	return err
}

func (svc Service) SavePoolTransactionEntry(message []byte) {
	var response rpc.NewTransactionSubscriptionResponse
	err := json.Unmarshal(message, &response)
	if err != nil {
		global.Logger.Error(err)
	}
	var result rpc.NewTransactionSubscriptionResult
	err = json.Unmarshal(json.RawMessage(response.Params.Result), &result)
	if err != nil {
		global.Logger.Error(err)
	}
	if result.Size != 0 {
		err = svc.CreateOrUpdatePoolTransactionEntry(result.Transaction, result.Fee, result.Cycles, result.Size)
		if err != nil {
			global.Logger.Errorf("PoolTransactionEntry creation error: %v", err)
		}
		global.Logger.Infof("receive: %s", string(message))
	}
}

func contains(supportedTopics []string, topic string) bool {
	for _, supportedTopic := range supportedTopics {
		if supportedTopic == topic {
			return true
		}
	}
	return false
}
