package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/AR1011/trade-engine/actors/tradeEngine"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/log"

	"github.com/google/uuid"
)

func main() {
	lh := log.NewHandler(os.Stdout, log.TextFormat, slog.LevelInfo)
	e := actor.NewEngine(actor.EngineOptLogger(log.NewLogger("[engine]", lh)))

	tradeEnginePID := e.Spawn(tradeEngine.NewTradeEngine(), "trade-engine")

	// create 5 trade orders
	// Expirary of 10s so after 10s the orders will be cancelled
	// the price watcher will be stopped due to inactivity

	for i := 0; i < 5; i++ {
		fmt.Println("Creating new trade order")
		o := &tradeEngine.TradeOrderRequest{
			TradeID:    uuid.New().String(),
			Token0:     "token0",
			Token1:     "token1",
			Chain:      "ETH",
			Wallet:     "random wallet",
			PrivateKey: "private key",
			// expire after 10 seconds
			Expires: time.Now().Add(time.Second * 10).UnixMilli(),
		}

		e.Send(tradeEnginePID, o)
	}

	// time.Sleep(time.Second * 5)
	// e.Send(tradeEnginePID, &tEngine.CancelOrderRequest{ID: trade1.TradeID})

	select {}

}
