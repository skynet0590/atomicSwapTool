// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package comms

import (
	"github.com/skynet0590/atomicSwapTool/dex"
)

// log is a logger that is initialized with no output filters. This means the
// package will not perform any logging by default until the caller requests it.
var log dex.Logger

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger dex.Logger) {
	log = logger
}
