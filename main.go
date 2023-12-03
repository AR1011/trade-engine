package main

import (
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

	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 2)
		o := &tradeEngine.TradeOrderRequest{
			TradeID: uuid.New().String(),
			Token0:  "0x000000000000000000",
			Token1:  "0x111111111111111111",
			Chain:   "ETH",
			Wallet:  "0x86bDd03525281214E2Ad874E616491D43c0233F2",
			Pk:      "289d095a1a421acb6498fecc656f5712d9aa95f63e8d9b321e162f28a2590f6f",
			// expire after 10 seconds
			Expires: time.Now().Add(time.Second * 10).UnixMilli(),
		}

		e.Send(tradeEnginePID, o)
	}

	// time.Sleep(time.Second * 5)
	// e.Send(tradeEnginePID, &tEngine.CancelOrderRequest{ID: trade1.TradeID})

	select {}

}
