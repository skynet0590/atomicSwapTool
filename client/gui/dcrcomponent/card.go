package dcrcomponent

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"image/color"
)

const (
	DefaultRadius = 10
)

type (
	Card struct {
		layout.Inset
		Color  color.NRGBA
		Radius CornerRadius
	}
)

func (c Card) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return c.Inset.Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						tr := float32(gtx.Px(unit.Dp(c.Radius.NE)))
						tl := float32(gtx.Px(unit.Dp(c.Radius.NW)))
						br := float32(gtx.Px(unit.Dp(c.Radius.SE)))
						bl := float32(gtx.Px(unit.Dp(c.Radius.SW)))
						clip.RRect{
							Rect: f32.Rectangle{Max: f32.Point{
								X: float32(gtx.Constraints.Min.X),
								Y: float32(gtx.Constraints.Min.Y),
							}},
							NE: tl, NW: tr, SE: br, SW: bl,
						}.Add(gtx.Ops)
						return fill(gtx, c.Color)
					}),
					layout.Stacked(w),
				)
			})
		}),
	)
	return dims
}

func fill(gtx layout.Context, col color.NRGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := image.Point{X: cs.Min.X, Y: cs.Min.Y}
	st := op.Save(gtx.Ops)
	track := image.Rectangle{
		Max: d,
	}
	clip.Rect(track).Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	st.Load()

	return layout.Dimensions{Size: d}
}
