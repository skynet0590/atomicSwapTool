package gui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/skynet0590/atomicSwapTool/client/gui/dcrcomponent"
	"github.com/skynet0590/atomicSwapTool/client/gui/validate/validators"
	"github.com/skynet0590/atomicSwapTool/client/gui/values"
)

type (
	registerPage struct {
		win          *Window
		registerForm registerForm
	}
	registerForm struct {
		win                 *Window
		passwordEditor      *dcrcomponent.Editor
		passwordAgainEditor *dcrcomponent.Editor
		submitButton        dcrcomponent.Button
	}
)

func (p *registerPage) Title() string {
	return "Register"
}

func (p *registerPage) HandlerEvent() {
	p.registerForm.handlerEvent()
}

func (p *registerPage) Layout(gtx C) D {
	p.HandlerEvent()
	return layout.Inset{
		Top:    unit.Dp(20),
		Right:  unit.Dp(20),
		Bottom: unit.Dp(20),
		Left:   unit.Dp(20),
	}.Layout(gtx, func(gtx layout.Context) D {
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
						return p.win.theme.H6("Set your app password. This password will protect your DEX account keys and connected wallets.").
							Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return p.win.theme.Line().Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return p.registerForm.Layout(gtx)
				}),
			)
		})
	})
}

func newRegisterForm(w *Window) registerForm {
	const passTxt = "Password"
	passwordEditor := w.theme.Editor(passTxt, validators.Required(passTxt))
	passwordEditor.Editor.Mask = '*'
	passwordAgainEditor := w.theme.Editor("Password Again", validators.MatchedInput(passTxt, passwordEditor.Editor))
	passwordAgainEditor.Editor.Mask = '*'
	return registerForm{
		win:                 w,
		passwordEditor:      &passwordEditor,
		passwordAgainEditor: &passwordAgainEditor,
		submitButton:        w.theme.Button("Submit"),
	}
}

func newRegisterPage(w *Window) *registerPage {
	page := registerPage{
		win:          w,
		registerForm: newRegisterForm(w),
	}
	return &page
}

func (f *registerForm) Layout(gtx C) D {
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
				return f.passwordEditor.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return f.passwordAgainEditor.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return f.submitButton.Layout(gtx)
			}),
		)
	})
}

func (f *registerForm) handlerEvent() {
	if f.submitButton.Button.Clicked() {
		fmt.Println("Clicked")
	}
}
