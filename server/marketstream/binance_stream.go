package marketstream

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/skynet0590/atomicSwapTool/dex"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

const binanceAddr = "stream.binance.com:9443"

type (
	TradingPair struct {
		Base  uint32
		Quote uint32
	}
	BinanceStream struct {
		Pairs []TradingPair
		lastTradings map[string]*AggTrade
	}
	binanceStreamData struct {
		Stream string          `json:"stream"`
		Data   json.RawMessage `json:"data"`
	}
	AggTrade struct {
		EventType        string `json:"e"` // Event type
		EventTime        uint64 `json:"E"` // Event time
		Symbol           string `json:"s"` // Symbol
		AggregateTradeID uint64 `json:"a"` // Aggregate trade ID
		Price            string `json:"p"` // Price
		Quantity         string `json:"q"` // Quantity
		FirstTradeID     uint64 `json:"f"` // First trade ID
		LastTradeID      uint64 `json:"l"` // Last trade ID
		TradeTime        uint64 `json:"T"` // Trade time
		IsBuyer          bool   `json:"m"` // Is the buyer the market maker?
		Ignore           bool   `json:"M"` // Ignore
	}
)

func (t *TradingPair) marketName() string {
	return dex.BipIDSymbol(t.Quote) + dex.BipIDSymbol(t.Base)
}

func NewBinanceStream(pairs ...TradingPair) (*BinanceStream, error) {
	if len(pairs) == 0 {
		return nil, fmt.Errorf("")
	}
	lastTradings := make(map[string]*AggTrade)
	for _, pair := range pairs {
		lastTradings[pair.marketName()] = nil
	}
	return &BinanceStream{
		Pairs: pairs,
		lastTradings: lastTradings,
	}, nil
}

func (s *BinanceStream) aggTradeQuery() string {
	var streamsQuery []string
	for _, pair := range s.Pairs {
		streamsQuery = append(streamsQuery, fmt.Sprintf("%s@aggTrade", pair.marketName()))
	}
	return strings.Join(streamsQuery, "/")
}

func (s *BinanceStream) Run(ctx context.Context) {
	// a := `{"stream":"dcrbtc@aggTrade","data":{"e":"aggTrade","E":1618132277929,"s":"DCRBTC","a":3603941,"p":"0.00324600","q":"0.03800000","f":5182754,"l":5182754,"T":1618132277928,"m":true,"M":true}}`

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: binanceAddr, Path: "/stream", RawQuery: fmt.Sprintf("streams=%s", s.aggTradeQuery())}
	log.Infof("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Info("read:", err)
				return
			}
			log.Infof("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			log.Info("Ticker: ", t.String())
		case <-interrupt:
			log.Info("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
