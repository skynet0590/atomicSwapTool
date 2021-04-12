package main

import (
	"context"
	"github.com/skynet0590/atomicSwapTool/dex"
	"github.com/skynet0590/atomicSwapTool/server/marketstream"
	"os"
)

func main() {
	marketstream.UseLogger(dex.NewLogger("MP", dex.LevelTrace, os.Stdout))
	stream, _ := marketstream.NewBinanceStream(marketstream.TradingPair{
		Base:  0,
		Quote: 42,
	})
	stream.Run(context.Background())
}
