package dcrcomponent

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"image"
	"image/color"
)

type (
	Theme struct {
		*material.Theme
		Clipboard Clipboard
		Color     Color
	}
	Color struct {
		Primary    color.NRGBA
		Secondary  color.NRGBA
		Text       color.NRGBA
		Hint       color.NRGBA
		Overlay    color.NRGBA
		InvText    color.NRGBA
		Success    color.NRGBA
		Danger     color.NRGBA
		Background color.NRGBA
		Surface    color.NRGBA
		Gray       color.NRGBA
		Black      color.NRGBA
		Orange     color.NRGBA
		LightGray  color.NRGBA
	}
)

var (
	DefaultColor = Color{
		Primary:    rgb(0x2970ff), // key blue
		Secondary:  rgb(0x091440), // dark blue
		Text:       rgb(0x091440),
		Hint:       rgb(0xbbbbbb),
		Overlay:    rgb(0x000000),
		InvText:    rgb(0xffffff),
		Success:    rgb(0x41bf53),
		Danger:     rgb(0xff0000),
		Background: argb(0x22444444),
		Surface:    rgb(0xffffff),
		Gray:       rgb(0x596D81),
		Black:      rgb(0x000000),
		Orange:     rgb(0xed6d47),
		LightGray:  rgb(0xc4cbd2),
	}
)

func NewTheme(fontCollection []text.FontFace, colorBoard Color, clipboard Clipboard) *Theme {
	theme := Theme{
		Theme:     material.NewTheme(fontCollection),
		Color:     colorBoard,
		Clipboard: clipboard,
	}
	return &theme
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func fillMax(gtx layout.Context, col color.NRGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Max.X, Y: cs.Max.Y}
	st := op.Save(gtx.Ops)
	track := image.Rectangle{
		Max: image.Point{X: d.X, Y: d.Y},
	}
	clip.Rect(track).Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	st.Load()
}

func FillMax(gtx layout.Context, col color.NRGBA) {
	fillMax(gtx, col)
}

func fill(gtx C, col color.NRGBA) D {
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

func Fill(gtx C, col color.NRGBA) D {
	return fill(gtx, col)
}
