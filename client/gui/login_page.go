package gui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/skynet0590/atomicSwapTool/client/gui/dcrcomponent"
	"github.com/skynet0590/atomicSwapTool/client/gui/values"
)

type (
	loginPage struct {
		win                 *Window
		passwordEditor      *dcrcomponent.Editor
		submitButton        dcrcomponent.Button
	}
)

func newLoginPage(win *Window) *loginPage {
	th := win.theme
	page := loginPage{
		win: win,
		passwordEditor:           th.EditorPassword("Wallet Password"),
	}

	return &page
}

func (p *loginPage) Layout(gtx C) D {
	p.handlerEvent()
	return layout.Inset{
		Top:    unit.Dp(20),
		Right:  unit.Dp(20),
		Bottom: unit.Dp(20),
		Left:   unit.Dp(20),
	}.Layout(gtx, func(gtx C) D {
		return p.win.theme.Card().Layout(gtx, func(gtx C) D {
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
						return p.win.theme.H6("Login").
							Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return p.win.theme.Line().Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
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
								return p.submitButton.Layout(gtx)
							}),
						)
					})
				}),
			)
		})
	})
}

func (p *loginPage) handlerEvent() {
	if p.passwordEditor.IsValid() {
		passTxt := p.passwordEditor.Text()
		_, err := p.win.core.Login([]byte(passTxt))

		if err == nil {
			p.win.loggedIn = true
		}else{
			p.win.Notify(fmt.Sprintf("Login failed: %v", err))
		}
	}
}

func (p *loginPage) Title() string {
	return "Login"
}
