package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	ckbRPC "github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/shaojunda/ckb-node-websocket-client/global"
	"github.com/shaojunda/ckb-node-websocket-client/internal/model"
	"github.com/shaojunda/ckb-node-websocket-client/internal/service"
	"github.com/shaojunda/ckb-node-websocket-client/pkg/logger"
	"github.com/shaojunda/ckb-node-websocket-client/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type doneCode struct {
	ExitCode int
}

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	err = setupCKBClient()
	if err != nil {
		log.Fatalf("init.setupCKBClient err: %v", err)
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: global.RPCSetting.WebSocketURL}
	global.Logger.Infof("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		global.Logger.Fatal("dial error: ", err)
	}
	defer c.Close()

	done := make(chan doneCode)

	svc := service.New(context.Background())
	err = svc.Subscribe(c, "new_transaction")
	if err != nil {
		global.Logger.Panic("write error: ", err)
		os.Exit(1)
	}

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				done <- doneCode{1}
				global.Logger.Fatal("node error: ", err)
				return
			}
			log.Printf("receive: %s", string(message))
			svc.SavePoolTransactionEntry(message)
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

func setupSetting() error {
	s, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("RPC", &global.RPCSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("SystemCodeHash", &global.SystemCodeHash)
	if err != nil {
		return err
	}

	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s%s", global.AppSetting.LogSavePath, global.AppSetting.LogFileName, global.AppSetting.LogFileExt),
		MaxSize:    600,
		MaxAge:     10,
		MaxBackups: 3,
		LocalTime:  true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

func setupCKBClient() error {
	client, err := ckbRPC.Dial(global.RPCSetting.URL)
	if err != nil {
		return err
	}
	global.CKBClient = client

	return nil
}
