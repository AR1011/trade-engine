package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/AR1011/trade-engine/actors/tEngine"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/log"

	"github.com/google/uuid"
)

func main() {
	lh := log.NewHandler(os.Stdout, log.TextFormat, slog.LevelDebug)
	e := actor.NewEngine(actor.Config{Logger: log.NewLogger("[engine]", lh)})

	tradeEnginePID := e.Spawn(tEngine.NewTradeEngine(), "trade-engine")

	order1 := &tEngine.TradeOrderRequest{
		TradeID: uuid.New().String(),
		Token0:  "0x000000000000000000",
		Token1:  "0x111111111111111111",
		Chain:   "ETH",
		Wallet:  "0x86bDd03525281214E2Ad874E616491D43c0233F2",
		Pk:      "289d095a1a421acb6498fecc656f5712d9aa95f63e8d9b321e162f28a2590f6f",
	}
	order2 := &tEngine.TradeOrderRequest{
		TradeID: uuid.New().String(),
		Token0:  "0x000000000000000000",
		Token1:  "0x111111111111111111",
		Chain:   "ETH",
		Wallet:  "0x86bDd03525281214E2Ad874E616491D43c0233F2",
		Pk:      "289d095a1a421acb6498fecc656f5712d9aa95f63e8d9b321e162f28a2590f6f",
	}

	e.Send(tradeEnginePID, order1)
	e.Send(tradeEnginePID, order2)

	// time.Sleep(time.Second * 5)
	// e.Send(tradeEnginePID, &tEngine.CancelOrderRequest{ID: trade1.TradeID})

	time.Sleep(time.Second * 50000)

}
