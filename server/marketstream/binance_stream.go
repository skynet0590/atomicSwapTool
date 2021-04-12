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
	"strconv"
	"strings"
	"sync"
)

const binanceAddr = "stream.binance.com:9443"

type (
	TradingPair struct {
		Base  uint32
		Quote uint32
	}
	BinanceStream struct {
		Pairs        []TradingPair
		mtx          sync.Mutex
		lastTradings map[string]TradingPrice
		url          url.URL
		callback     func(symbol string, price TradingPrice)
	}
	binanceStreamData struct {
		Stream string          `json:"stream"`
		Data   json.RawMessage `json:"data"`
	}
	TradingPrice struct {
		Symbol    string
		Price     float64
		TradeTime uint64
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
	Trading interface {
		GetSymbol() string
		GetPrice() float64
		GetTradeTime() uint64
	}
)

func (t *TradingPair) marketName() string {
	return dex.BipIDSymbol(t.Quote) + dex.BipIDSymbol(t.Base)
}

func NewBinanceStream(callback func(symbol string, price TradingPrice), pairs ...TradingPair) (*BinanceStream, error) {
	if len(pairs) == 0 {
		return nil, fmt.Errorf("")
	}
	stream := BinanceStream{
		Pairs:    pairs,
		callback: callback,
	}
	stream.init()
	return &stream, nil
}

func (s *BinanceStream) init() {
	var streamsQuery []string
	lastTradings := make(map[string]TradingPrice)
	for _, pair := range s.Pairs {
		lastTradings[pair.marketName()] = TradingPrice{}
		streamsQuery = append(streamsQuery, fmt.Sprintf("%s@aggTrade", pair.marketName()))
	}
	s.lastTradings = lastTradings
	s.url = url.URL{
		Scheme:   "wss",
		Host:     binanceAddr,
		Path:     "/stream",
		RawQuery: fmt.Sprintf("streams=%s", strings.Join(streamsQuery, "/")),
	}
}

func (s *BinanceStream) Run(ctx context.Context) {
	if err := s.connect(); err != nil {
		log.Error("Connect to Binance failed: ", err)
	}
}

func (s *BinanceStream) parseMessageData(msg []byte) error {
	streamData := binanceStreamData{}
	err := json.Unmarshal(msg, &streamData)
	if err != nil {
		return err
	}
	ok, symbol := streamData.isAggTrade()
	if !ok {
		return fmt.Errorf("Message data is not supported yet.")
	}
	aggTrade := AggTrade{}
	err = json.Unmarshal(streamData.Data, &aggTrade)
	if err != nil {
		return err
	}
	s.updateTradingInfo(symbol, &aggTrade)
	return err
}

func (s *BinanceStream) updateTradingInfo(symbol string, trade Trading) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	trading := TradingPrice{
		Symbol:    trade.GetSymbol(),
		Price:     trade.GetPrice(),
		TradeTime: trade.GetTradeTime(),
	}
	if s.callback != nil {
		s.callback(symbol, trading)
	}
	s.lastTradings[symbol] = trading
}

func (s *BinanceStream) connect() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	log.Infof("Connecting to %s", s.url.String())
	c, _, err := websocket.DefaultDialer.Dial(s.url.String(), nil)
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
				log.Info("Closed connection from Binance server:", err)
				return
			}
			err = s.parseMessageData(message)
			if err != nil {
				log.Error("Parse message data:", err)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			log.Info("Try to reconnect")
			return s.connect()
		case <-interrupt:
			log.Info("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Info("Closing web socket connection:", err)
				return err
			}
			return nil
		}
	}
}

func (s *BinanceStream) GetTrade(pair TradingPair) TradingPrice {
	return s.GetTradeBySymbol(pair.marketName())
}

func (s *BinanceStream) GetTradeBySymbol(symbol string) TradingPrice {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	trade, _ := s.lastTradings[symbol]
	return trade
}

func (d *binanceStreamData) isAggTrade() (ok bool, name string) {
	info := strings.Split(d.Stream, "@")
	if len(info) != 2 {
		return false, ""
	}
	if info[1] == "aggTrade" {
		return true, info[0]
	}
	return false, ""
}

func (t *AggTrade) GetSymbol() string {
	return t.Symbol
}

func (t *AggTrade) GetPrice() float64 {
	price, _ := strconv.ParseFloat(t.Price, 64)
	return price
}

func (t *AggTrade) GetTradeTime() uint64 {
	return t.TradeTime
}
