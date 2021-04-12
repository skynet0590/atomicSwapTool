package main

import (
	"context"
	"github.com/skynet0590/atomicSwapTool/dex"
	"github.com/skynet0590/atomicSwapTool/server/marketstream"
	"os"
)

func main() {
	log := dex.NewLogger("MP", dex.LevelTrace, os.Stdout)
	marketstream.UseLogger(log)
	stream, _ := marketstream.NewBinanceStream(func(symbol string, price marketstream.TradingPrice) {
		log.Info(symbol, price)
	},marketstream.TradingPair{
		Base:  0,
		Quote: 42,
	})
	stream.Run(context.Background())
}
