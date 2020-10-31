package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/shaojunda/ckb-node-websocket-client/global"
	"github.com/shaojunda/ckb-node-websocket-client/internal/model"
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
