// +build harness

package livetest

import (
	"fmt"
	"testing"

	"github.com/skynet0590/atomicSwapTool/client/asset/btc"
	"github.com/skynet0590/atomicSwapTool/dex"
	dexbtc "github.com/skynet0590/atomicSwapTool/dex/networks/btc"
)

const (
	alphaAddress = "bcrt1qy7agjj62epx0ydnqskgwlcfwu52xjtpj36hr0d"
)

var (
	tBTC = &dex.Asset{
		ID:           0,
		Symbol:       "btc",
		SwapSize:     dexbtc.InitTxSizeSegwit,
		SwapSizeBase: dexbtc.InitTxSizeBaseSegwit,
		MaxFeeRate:   10,
		LotSize:      1e6,
		RateStep:     10,
		SwapConf:     1,
	}
)

func TestWallet(t *testing.T) {
	fmt.Println("////////// WITHOUT SPLIT FUNDING TRANSACTIONS //////////")
	Run(t, btc.NewWallet, alphaAddress, tBTC, false)
	fmt.Println("////////// WITH SPLIT FUNDING TRANSACTIONS //////////")
	Run(t, btc.NewWallet, alphaAddress, tBTC, true)
}
