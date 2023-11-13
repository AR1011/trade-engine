package executor

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/AR1011/trade-engine/actors/price"
	"github.com/AR1011/trade-engine/utils"
	"github.com/anthdm/hollywood/actor"
)

type TradeExecutorKillRequest struct {
	ID string
}

// message to get trade info
type TradeInfoRequest struct{}

// response message for trade info
type TradeInfoResponse struct {
	// info regarding the current position
	// eg price, pnl, etc
	foo   int
	bar   int
	price float64 // will be decimal
}

type ExecutorOptions struct {
	PriceWatcherPID *actor.PID
	TradeID         string //uuid string
	Ticker          string
	Token0          string
	Token1          string
	Chain           string
	Wallet          string
	Pk              string
}

type tradeExecutor struct {
	id              string // uuid string
	actorEngine     *actor.Engine
	tradeEnginePID  *actor.PID
	priceWatcherPID *actor.PID
	status          string
	ticker          string
	token0          string
	token1          string
	chain           string
	wallet          string
	pk              string
	price           float64 // will be decimal
	active          bool
	// ... will contain more
	// will also contain runtime vars
}

func (te *tradeExecutor) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		slog.Info(utils.TExc+utils.PadP("Init Trade Executor Actor"), "id", te.id, "wallet", te.wallet)

		// set flag for goroutine
		te.active = true

		// start the trade process
		go te.init(c)

		te.actorEngine = c.Engine()
		te.tradeEnginePID = c.GetPID("trade-engine")

	case actor.Stopped:
		slog.Info(utils.TExc+utils.PadP("Stopped Trade Executor Actor"), "id", te.id, "wallet", te.wallet)
		te.active = false

	case TradeInfoRequest:
		slog.Info(utils.TExc+utils.PadP("Got TradeInfoRequest"), "id", te.id, "wallet", te.wallet)
		te.tradeInfo(c)

	default:
		_ = msg

	}
}

func (te *tradeExecutor) init(c *actor.Context) {
	var i int

	for {

		// check flag. Will be false if actor is killed
		if !te.active {
			return
		}

		// for demo / testing, after 5 iterations, kill the actor
		if i > 5 {
			te.Finished()
			return
		}

		time.Sleep(time.Second * 5)

		if (te.priceWatcherPID == nil) || (te.priceWatcherPID == &actor.PID{}) {
			slog.Error(utils.TExc + utils.PadP("priceWatcherPID is <nil>"))
			return
		}

		// get the price from the price actor, 1s timeout
		response := c.Request(te.priceWatcherPID, price.FetchPriceRequest{}, time.Second)

		// wait for result
		result, err := response.Result()
		if err != nil {
			// fuck!!
			slog.Error(utils.TExc+utils.PadP("Error getting price response"), "error", err.Error())
			return
		}

		switch r := result.(type) {
		case *price.FetchPriceResponse:
			slog.Info(utils.TExc+utils.PadP("Got Price Response"), "price", r.Price)
		default:
			slog.Warn(utils.TExc+utils.PadP("Got Invalid Type from priceWatcher"), "type", reflect.TypeOf(r))

		}

		i++
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

	slog.Info(utils.TExc+utils.PadP("Finished"), "id", te.id)
	// set the flag to flase so goroutine returns
	te.active = false

	// send kill request to the trade engine so it can remove it from maps
	// and poision the actor

	// make sure tradeEnginePID and actorEngine are safe
	if te.tradeEnginePID == nil {
		slog.Error(utils.TExc + utils.PadP("tradeEnginePID is <nil>"))
	}

	if te.actorEngine == nil {
		slog.Error(utils.TExc + utils.PadP("actorEngine is <nil>"))
	}

	te.actorEngine.Send(te.tradeEnginePID, &TradeExecutorKillRequest{ID: te.id})
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
		}
	}
}
