package price

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/AR1011/trade-engine/utils"
	"github.com/anthdm/hollywood/actor"
)

type PriceOptions struct {
	Ticker string
	Token0 string
	Token1 string
	Chain  string
}

type PriceWatcherKillRequest struct {
	Ticker string
}

type FetchPriceRequest struct{}

type FetchPriceResponse struct {
	Iat   int64
	Price float64
}

type priceWatcher struct {
	actorEngine    *actor.Engine
	tradeEnginePID *actor.PID
	ticker         string
	token0         string
	token1         string
	chain          string
	lastPrice      float64 // will use decimal in real
	updatedAt      int64
	lastCall       int64
	// will contain more stuff
}

func (pw *priceWatcher) Receive(c *actor.Context) {

	switch msg := c.Message().(type) {
	case actor.Started:

	case actor.Initialized:
		slog.Info(utils.PWat+utils.PadG("Init Price Actor"), "ticker", pw.ticker)

		// start updating the price
		go pw.init()

		pw.actorEngine = c.Engine()
		pw.tradeEnginePID = c.GetPID("trade-engine")
		pw.lastCall = time.Now().UnixMilli()

	case actor.Stopped:
		slog.Info(utils.PWat+utils.PadG("Stopped Price Actor"), "ticker", pw.ticker)

	case FetchPriceRequest:
		slog.Info(utils.PWat+utils.PadG("Fetching Price Request"), "ticker", pw.ticker)

		// update last called time
		pw.lastCall = time.Now().UnixMilli()

		// respond with the lastest price
		c.Respond(&FetchPriceResponse{
			Iat:   time.Now().UnixMilli(),
			Price: pw.lastPrice,
		})

	default:
		slog.Warn(utils.PWat+utils.PadG("Got Invalid Message Type"), "ticker", pw.ticker, "type", reflect.TypeOf(msg))

		_ = msg
	}
}

func (pw *priceWatcher) init() {
	// mimic getting price every 2 seconds
	for {
		// check if the last call was more than 30 seconds ago
		if pw.lastCall < time.Now().UnixMilli()-(time.Second.Milliseconds()*30) {
			// if no call in 30 seconds => kill itself
			pw.Kill()
			return // stops goroutine
		}

		time.Sleep(time.Second * 2)
		pw.lastPrice++
		pw.updatedAt = time.Now().UnixMilli()

	}
}

func (pw *priceWatcher) Kill() {
	// send kill request to the trade engine so it can remove it from maps
	// and poision the actor

	// make sure tradeEnginePID and actorEngine are safe
	if pw.tradeEnginePID == nil {
		slog.Error(utils.PWat+utils.PadG("tradeEnginePID is <nil>"), "ticker", pw.ticker)
	}

	if pw.actorEngine == nil {
		slog.Error(utils.PWat+utils.PadG("actorEngine is <nil>"), "ticker", pw.ticker)
	}

	pw.actorEngine.Send(pw.tradeEnginePID, &PriceWatcherKillRequest{Ticker: pw.ticker})
}

func NewPriceActor(opts PriceOptions) actor.Producer {
	return func() actor.Receiver {
		return &priceWatcher{
			ticker: opts.Ticker,
			token0: opts.Token0,
			token1: opts.Token1,
			chain:  opts.Chain,
		}
	}
}
