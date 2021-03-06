package dcrcomponent

import "gioui.org/layout"

type (
	ReadClipboard struct{}
	C = layout.Context
	D = layout.Dimensions
	CornerRadius struct {
		NE float32
		NW float32
		SE float32
		SW float32
	}
)
