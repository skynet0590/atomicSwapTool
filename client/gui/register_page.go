package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/skynet0590/atomicSwapTool/client/gui/dcrcomponent"
	"github.com/skynet0590/atomicSwapTool/client/gui/values"
)

type (
	registerPage struct {
		win *Window
		passwordEditor dcrcomponent.Editor
		passwordAgainEditor dcrcomponent.Editor
	}
)

func (p *registerPage) Title() string {
	return "Register page"
}

func (p *registerPage) Layout(gtx C) D {
	return layout.Inset{
		Top:    unit.Dp(20),
		Right:  unit.Dp(20),
		Bottom: unit.Dp(20),
		Left:   unit.Dp(20),
	}.Layout(gtx, func(gtx layout.Context) D {
		return p.win.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding15,
						Right:  values.MarginPadding10,
						Bottom: values.MarginPadding10,
						Left:   values.MarginPadding10,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return p.win.H6("Set your app password. This password will protect your DEX account keys and connected wallets.").
							Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return p.win.Line().Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return p.passwordComponent(gtx)
				}),
			)
		})
	})
}

func (p *registerPage) passwordComponent(gtx C) D {
	return layout.Inset{
		Top:    values.MarginPadding15,
		Right:  values.MarginPadding10,
		Bottom: values.MarginPadding10,
		Left:   values.MarginPadding10,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return p.passwordEditor.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return p.passwordAgainEditor.Layout(gtx)
			}),
		)

	})
}

func newRegisterPage(w *Window) *registerPage {
	passwordEditor := w.Editor( "Password")
	passwordEditor.Editor.SingleLine = true
	passwordAgainEditor := w.Editor("Password Again")
	passwordAgainEditor.Editor.SingleLine = true
	return &registerPage{
		win:w,
		passwordEditor: passwordEditor,
		passwordAgainEditor: passwordAgainEditor,
	}
}
