package gui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image/color"
)

type (
	overviewPage struct {
		win *Window
	}
)

func (p *overviewPage) Title() string {
	return "Overview page"
}

func (p *overviewPage) Layout(gtx C) D {
	th := p.win.theme.Theme
	return layout.Inset{
		Top:    unit.Dp(20),
		Right:  unit.Dp(20),
		Bottom: unit.Dp(20),
		Left:   unit.Dp(20),
	}.Layout(gtx, func(gtx layout.Context) D {
		l := material.H5(th, p.Title())
		maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
		l.Color = maroon
		l.Alignment = text.Start
		return l.Layout(gtx)
	})
}

func newOverviewPage(w *Window) *overviewPage {
	return &overviewPage{
		win:w,
	}
}
