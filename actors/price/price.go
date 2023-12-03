package price

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/shopspring/decimal"
)

type PriceOptions struct {
	Ticker string
	Token0 string
	Token1 string
	Chain  string
}

type FetchPriceRequest struct{}

type FetchPriceResponse struct {
	Iat   int64
	Price decimal.Decimal
}

type priceWatcher struct {
	actorEngine *actor.Engine
	PID         *actor.PID
	ticker      string
	token0      string
	token1      string
	chain       string
	lastPrice   decimal.Decimal
	updatedAt   int64
	lastCall    int64
	callCount   uint64
}

func (pw *priceWatcher) Receive(c *actor.Context) {

	switch msg := c.Message().(type) {
	case actor.Started:
		slog.Info("Started Price Actor", "ticker", pw.ticker)

		pw.actorEngine = c.Engine()
		pw.lastCall = time.Now().UnixMilli()
		pw.PID = c.PID()

		// start updating the price
		go pw.start()

	case actor.Stopped:
		slog.Info("Stopped Price Actor", "ticker", pw.ticker)

	case FetchPriceRequest:
		slog.Info("Fetching Price Request", "ticker", pw.ticker)

		// update last called time
		pw.lastCall = time.Now().UnixMilli()

		// increment call count
		pw.callCount++

		// respond with the lastest price
		c.Respond(&FetchPriceResponse{
			Iat:   time.Now().UnixMilli(),
			Price: pw.lastPrice,
		})

	default:
		slog.Warn("Got Invalid Message Type", "ticker", pw.ticker, "type", reflect.TypeOf(msg))
		_ = msg
	}
}

func (pw *priceWatcher) start() {
	pw.lastPrice = decimal.NewFromInt(0)

	// mimic getting price every 2 seconds
	for {
		// check if the last call was more than 10 seconds ago
		if pw.lastCall < time.Now().UnixMilli()-(time.Second.Milliseconds()*10) {
			slog.Warn("Inactivity: Killing Price Watcher", "ticker", pw.ticker, "callCount", pw.callCount)

			// if no call in 10 seconds => kill itself
			pw.Kill()
			return // stops goroutine
		}

		time.Sleep(time.Millisecond * 2)
		pw.lastPrice = pw.lastPrice.Add(decimal.NewFromFloat(1))
		pw.updatedAt = time.Now().UnixMilli()

	}
}

func (pw *priceWatcher) Kill() {
	// send kill request to the trade engine so it can remove it from maps
	// and poision the actor

	if pw.actorEngine == nil {
		slog.Error("actorEngine is <nil>", "ticker", pw.ticker)
	}
	pw.actorEngine.Poison(pw.PID)
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
