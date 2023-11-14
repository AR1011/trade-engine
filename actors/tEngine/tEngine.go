package tEngine

import (
	"fmt"

	"github.com/AR1011/trade-engine/actors/executor"
	"github.com/AR1011/trade-engine/actors/price"
	"github.com/AR1011/trade-engine/logger"
	"github.com/anthdm/hollywood/actor"
)

type tradeEngine struct {
	logger logger.Logger
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
		t.logger.Info("Init Trade Engine")
		_ = msg

	case *TradeOrderRequest:
		// got new trade order, create the executor
		t.logger.Info("Got New Trade Order", "id", msg.TradeID, "wallet", msg.Wallet)
		t.spawnExecutor(msg, c)

	}
}

func (t *tradeEngine) spawnExecutor(msg *TradeOrderRequest, c *actor.Context) {
	// make sure is price stream for the pair
	pricePID := t.ensurePriceStream(msg, c)

	// spawn the executor
	options := &executor.ExecutorOptions{
		PriceWatcherPID: pricePID,
		TradeID:         msg.TradeID,
		Ticker:          toTicker(msg.Token0, msg.Token1, msg.Chain),
		Token0:          msg.Token0,
		Token1:          msg.Token1,
		Chain:           msg.Chain,
		Wallet:          msg.Wallet,
		Pk:              msg.Pk,
	}

	// spawn the actor
	c.SpawnChild(executor.NewExecutorActor(options), msg.TradeID)

}

func (t *tradeEngine) ensurePriceStream(order *TradeOrderRequest, c *actor.Context) *actor.PID {
	ticker := toTicker(order.Token0, order.Token1, order.Chain)

	pid := c.Child("trade-engine/" + ticker)
	if pid != nil {
		return pid
	}

	options := price.PriceOptions{
		Ticker: ticker,
		Token0: order.Token0,
		Token1: order.Token1,
		Chain:  order.Chain,
	}

	// spawn the actor
	pid = c.SpawnChild(price.NewPriceActor(options), ticker)
	return pid
}

func NewTradeEngine() actor.Producer {
	return func() actor.Receiver {
		return &tradeEngine{
			logger: logger.NewLogger(
				logger.TEng,
				logger.DBLUE,
				logger.WithToStdoutWriter(),
				logger.WithToFileWriter("./logs/trade-engine.log", logger.JSON),
			),
		}
	}
}

func toTicker(token0, token1, chain string) string {
	return fmt.Sprintf("%s-%s-%s", token0, token1, chain)
}
