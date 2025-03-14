package executor

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/AR1011/trade-engine/actors/price"
	"github.com/anthdm/hollywood/actor"
	"github.com/shopspring/decimal"
)

// message to get trade info
type TradeInfoRequest struct{}

// response message for trade info
type TradeInfoResponse struct {
	// info regarding the current position
	// eg price, pnl, etc
	foo   int
	bar   int
	price decimal.Decimal
}

type ExecutorOptions struct {
	PriceWatcherPID *actor.PID
	TradeID         string
	Ticker          string
	Token0          string
	Token1          string
	Chain           string
	Wallet          string
	Pk              string
	Expires         int64
}

type tradeExecutor struct {
	id              string
	actorEngine     *actor.Engine
	PID             *actor.PID
	priceWatcherPID *actor.PID
	Expires         int64
	status          string
	ticker          string
	token0          string
	token1          string
	chain           string
	wallet          string
	pk              string
	price           decimal.Decimal
	active          bool
}

func (te *tradeExecutor) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		slog.Info("Started Trade Executor Actor", "id", te.id, "wallet", te.wallet)

		// set flag for goroutine
		te.active = true

		te.actorEngine = c.Engine()
		te.PID = c.PID()

		// start the trade process
		go te.start(c)

	case actor.Stopped:
		slog.Info("Stopped Trade Executor Actor", "id", te.id, "wallet", te.wallet)
		te.active = false

	case TradeInfoRequest:
		slog.Info("Got TradeInfoRequest", "id", te.id, "wallet", te.wallet)
		te.tradeInfo(c)

	default:
		_ = msg

	}
}

func (te *tradeExecutor) start(c *actor.Context) {
	// example of a long running process
	for {
		// check flag. Will be false if actor is killed
		if !te.active {
			return
		}

		if time.Now().UnixMilli() > te.Expires {
			slog.Warn("Trade Expired", "id", te.id, "wallet", te.wallet)
			te.Finished()
			return
		}

		// refresh price every 2s
		time.Sleep(time.Second * 2)

		if (te.priceWatcherPID == nil) || (te.priceWatcherPID == &actor.PID{}) {
			slog.Error("priceWatcherPID is <nil>")
			return
		}

		// get the price from the price actor, 2s timeout
		response := c.Request(te.priceWatcherPID, price.FetchPriceRequest{}, time.Second*2)

		// wait for result
		result, err := response.Result()
		if err != nil {
			slog.Error("Error getting price response", "error", err.Error())
			return
		}

		switch r := result.(type) {
		case *price.FetchPriceResponse:
			slog.Info("Got Price Response", "price", r.Price.StringFixed(18))
		default:
			slog.Warn("Got Invalid Type from priceWatcher", "type", reflect.TypeOf(r))

		}
	}
}

func (te *tradeExecutor) tradeInfo(c *actor.Context) {
	c.Respond(&TradeInfoResponse{
		foo:   100,
		bar:   100,
		price: te.price,
	})
}

func (te *tradeExecutor) Finished() {
	// set the flag to flase so goroutine terminates
	te.active = false

	// make sure actorEngine is safe
	if te.actorEngine == nil {
		slog.Error("actorEngine is <nil>")

	}
	te.actorEngine.Poison(te.PID)
}

func NewExecutorActor(opts *ExecutorOptions) actor.Producer {
	return func() actor.Receiver {
		return &tradeExecutor{
			id:              opts.TradeID,
			ticker:          opts.Ticker,
			token0:          opts.Token0,
			token1:          opts.Token1,
			chain:           opts.Chain,
			wallet:          opts.Wallet,
			pk:              opts.Pk,
			priceWatcherPID: opts.PriceWatcherPID,
			Expires:         opts.Expires,
			status:          "pending",
		}
	}
}
