package gui

import "github.com/skynet0590/atomicSwapTool/client/gui/dcrcomponent"

func (t *Window) Card() dcrcomponent.Card {
	return dcrcomponent.Card{
		Color: t.Color.Surface,
		Radius: dcrcomponent.CornerRadius{
			NE: dcrcomponent.DefaultRadius,
			SE: dcrcomponent.DefaultRadius,
			NW: dcrcomponent.DefaultRadius,
			SW: dcrcomponent.DefaultRadius,
		},
	}
}
