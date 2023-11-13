package tEngine

import (
	"fmt"
	"log/slog"

	"github.com/AR1011/trade-engine/actors/executor"
	"github.com/AR1011/trade-engine/actors/price"
	"github.com/anthdm/hollywood/actor"
)

type tradeEngine struct {
	pricePIDs    map[string]*actor.PID
	executorPIDs map[string]*actor.PID
}

type TradeOrderRequest struct {
	//will contain more
	TradeID string //uuid string
	Token0  string
	Token1  string
	Chain   string
	Wallet  string
	Pk      string
}

type CancelOrderRequest struct {
	ID string
}

func (t *tradeEngine) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Stopped:
		// should propogate to children and kill them

	case actor.Initialized:
		slog.Info("[TRADE ENGINE] Init Trade Engine Actor")

		_ = msg

	case *TradeOrderRequest:
		// got new trade order, create the executor
		slog.Info("[TRADE ENGINE] Got New TradeOrderRequest")
		t.spawnExecutor(msg, c)

	case *price.PriceWatcherKillRequest:
		slog.Info("[TRADE ENGINE] Killing Inactive Price Watcher")
		t.killPriceWatcher(msg, c)

	case *CancelOrderRequest:
		slog.Info("[TRADE ENGINE] Cancelling Order")
		t.killTradeExecutor(msg, c)

	}

}

func (t *tradeEngine) spawnExecutor(msg *TradeOrderRequest, c *actor.Context) {
	// make sure is price stream for the pair
	pricePID := t.ensurePriceStream(msg, c)

	// spawn the executor
	options := &executor.ExecutorOptions{
		PriceWatcherPID: pricePID,
		TradeID:         msg.TradeID,
		Ticker:          fmt.Sprintf("%s/%s/%s", msg.Token0, msg.Token1, msg.Chain),
		Token0:          msg.Token0,
		Token1:          msg.Token1,
		Chain:           msg.Chain,
		Wallet:          msg.Wallet,
		Pk:              msg.Pk,
	}

	// spawn the actor
	pid := c.SpawnChild(executor.NewExecutorActor(options), msg.TradeID)

	// store the pid
	t.executorPIDs[msg.TradeID] = pid

}

func (t *tradeEngine) ensurePriceStream(order *TradeOrderRequest, c *actor.Context) *actor.PID {
	ticker := fmt.Sprintf("%s/%s/%s", order.Token0, order.Token1, order.Chain)

	// check if there is an existing PID for the same ticker
	if pid, found := t.executorPIDs[ticker]; found {
		slog.Info("[TRADE ENGINE] Found Existing Price Watcher", "ticker", ticker)
		return pid
	} else {
		// if not then create new price watcher
		options := price.PriceOptions{
			Ticker: ticker,
			Token0: order.Token0,
			Token1: order.Token1,
			Chain:  order.Chain,
		}

		// spawn the actor
		pid = c.SpawnChild(price.NewPriceActor(options), ticker)

		// store the pid
		t.executorPIDs[ticker] = pid
		slog.Info("[TRADE ENGINE] Spawned New Price Watcher", "ticker", ticker)
		return pid
	}
}

func (t *tradeEngine) killPriceWatcher(req *price.PriceWatcherKillRequest, c *actor.Context) {
	// check if pid map has the ticker
	pid, ok := t.pricePIDs[req.Ticker]
	if !ok {
		// if not then return
		return
	}

	// kill the actor
	c.Engine().Poison(pid)
}

func (t *tradeEngine) killTradeExecutor(req *CancelOrderRequest, c *actor.Context) {
	// check if pid map has the ticker
	pid, ok := t.executorPIDs[req.ID]
	if !ok {
		// if not then return
		return
	}

	// kill the actor
	c.Engine().Poison(pid)
}

func NewTradeEngine() actor.Producer {
	return func() actor.Receiver {
		return &tradeEngine{
			pricePIDs:    make(map[string]*actor.PID),
			executorPIDs: make(map[string]*actor.PID),
		}
	}
}
