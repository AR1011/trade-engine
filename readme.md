# trade-engine example demo

<p align="center">
  <img src="assets/diagram.png" alt="diagram" width="100%">
</p>

## There are 3 actor types:

1. Trade Engine: There will just be one, that will be used to create and manage the Price Watcher and Trade Executor Actors

2. Price Watcher: Actor that will get the price for a given ticker. There will only be one actor for each ticker to save on resources.

   - The Price Watcher starts a go routine on `actor.Initialised` message. The go routine refreshes the price every x seconds.
   - If after 10 seconds no one has requested the latest price, the go routine is stopped and the actor is poisoned.

3. Trade Executor: Actor that will execute trades. There will be one actor for each trade.
   - Similar to Price Watcher, the Trade Executor will create a go routine on `actor.Initialised` message.
   - The go routine will periodically get the latest price from the Price Watcher by sending a `FetchPriceRequest` message. It will monitor the price and execute the trade when the price is right.
   - If the trade is canceled, a flag will be set that causes the go routine to stop, and the actor will be poisoned by the Trade Engine.

All the trade states will be stored in DB, so that if the Trade Engine crashes, it can recover the state of all the trades.
