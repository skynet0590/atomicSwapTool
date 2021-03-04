package gui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image/color"
)

type (
	registerPage struct {
		win *Window
	}
)

func (p *registerPage) Title() string {
	return "Register page"
}

func (p *registerPage) Layout(gtx C) D {
	th := p.win.theme
	return layout.Inset{
		Top:    unit.Dp(20),
		Right:  unit.Dp(20),
		Bottom: unit.Dp(20),
		Left:   unit.Dp(20),
	}.Layout(gtx, func(gtx layout.Context) D {
		return p.win.Card().Layout(gtx, func(gtx C) D {
			l := material.H5(th, p.Title())
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			l.Color = maroon
			l.Alignment = text.Start
			return l.Layout(gtx)
		})
	})
}

func newRegisterPage(w *Window) *registerPage {
	return &registerPage{
		win:w,
	}
}
