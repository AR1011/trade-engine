package price

import (
	"reflect"
	"time"

	"github.com/AR1011/trade-engine/logger"
	"github.com/anthdm/hollywood/actor"
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
	Price float64
}

type priceWatcher struct {
	actorEngine *actor.Engine
	PID         *actor.PID
	ticker      string
	token0      string
	token1      string
	chain       string
	lastPrice   float64 // will use decimal in real
	updatedAt   int64
	lastCall    int64
	callCount   uint64
	logger      logger.Logger
	// will contain more stuff
}

func (pw *priceWatcher) Receive(c *actor.Context) {

	switch msg := c.Message().(type) {
	case actor.Started:

	case actor.Initialized:
		pw.logger.Info("Init Price Actor", "ticker", pw.ticker)

		pw.actorEngine = c.Engine()
		pw.lastCall = time.Now().UnixMilli()
		pw.PID = c.PID()

		// start updating the price
		go pw.init()

	case actor.Stopped:
		pw.logger.Info("Stopped Price Actor", "ticker", pw.ticker)

	case FetchPriceRequest:
		pw.logger.Info("Fetching Price Request", "ticker", pw.ticker)

		// update last called time
		pw.lastCall = time.Now().UnixMilli()
		pw.callCount++

		// respond with the lastest price
		c.Respond(&FetchPriceResponse{
			Iat:   time.Now().UnixMilli(),
			Price: pw.lastPrice,
		})

	default:
		pw.logger.Warn("Got Invalid Message Type", "ticker", pw.ticker, "type", reflect.TypeOf(msg))

		_ = msg
	}
}

func (pw *priceWatcher) init() {
	// mimic getting price every 2 seconds
	for {
		// check if the last call was more than 10 seconds ago
		if pw.lastCall < time.Now().UnixMilli()-(time.Second.Milliseconds()*10) {
			pw.logger.Warn("Inactivity: Killing Price Watcher", "ticker", pw.ticker, "callCount", pw.callCount)

			// if no call in 10 seconds => kill itself
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

	if pw.actorEngine == nil {
		pw.logger.Error("actorEngine is <nil>", "ticker", pw.ticker)
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
			logger: logger.NewLogger(
				logger.PriceWatcher,
				logger.ColorDarkPurple,
				logger.LevelInfo,
				logger.WithToStdoutWriter(),
				logger.WithToFileWriter("./logs/trade-engine.log", logger.JsonFormat),
			),
		}
	}
}
