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
		win           *Window
		registerForm  *registerForm
		dcrWalletForm *dcrWalletSettingForm
		dexServerForm *dexServerForm
		confirmForm   *confirmForm
		step          registerStep
	}
	registerForm struct {
		win                 *Window
		parents             *registerPage
		passwordEditor      *dcrcomponent.Editor
		passwordAgainEditor *dcrcomponent.Editor
		submitButton        dcrcomponent.Button
		callback            func()
	}
	registerStep uint
)

const (
	registerSetPassword registerStep = 1 << iota
	registerDCRWallet
	registerAddDEXform
	registerConfirmForm
)

func (p *registerPage) Title() string {
	return "Register"
}

func (p *registerPage) handlerEvent() {
	p.registerForm.handlerEvent()
}

func (p *registerPage) Layout(gtx C) D {
	p.handlerEvent()
	return layout.Inset{
		Top:    unit.Dp(20),
		Right:  unit.Dp(20),
		Bottom: unit.Dp(20),
		Left:   unit.Dp(20),
	}.Layout(gtx, func(gtx layout.Context) D {
		return p.win.theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			switch p.step {
			case registerSetPassword:
				return p.registerForm.Layout(gtx)
			case registerDCRWallet:
				return p.dcrWalletForm.Layout(gtx)
			case registerAddDEXform:
				return p.dexServerForm.Layout(gtx)
			default:
				p.win.ChangePage(overview)
				return D{}
			}
		})
	})
}

func newRegisterForm(w *Window, callback func()) *registerForm {
	const passTxt = "Password"
	passwordEditor := w.theme.EditorPassword(passTxt, validators.Required(passTxt))
	passwordAgainEditor := w.theme.EditorPassword("Password Again", validators.MatchedInput(passTxt, passwordEditor.Editor))
	return &registerForm{
		win:                 w,
		passwordEditor:      passwordEditor,
		passwordAgainEditor: passwordAgainEditor,
		submitButton:        w.theme.Button("Submit"),
		callback:            callback,
	}
}

func newRegisterPage(w *Window) *registerPage {
	page := &registerPage{
		win:  w,
		step: registerSetPassword,
	}
	page.registerForm = newRegisterForm(w, func() {
		page.step = registerDCRWallet
	})
	page.dcrWalletForm = newDcrWalletSetting(w, func() {
		page.step = registerAddDEXform
	})
	page.dexServerForm = newDEXServerForm(w, func() {
		w.ChangePage(overview)
	})
	return page
}

func (f *registerForm) Layout(gtx C) D {
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
				return f.win.theme.H6("Set your app password. This password will protect your DEX account keys and connected wallets.").
					Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return f.win.theme.Line().Layout(gtx)
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
		}),
	)
}

func (f *registerForm) handlerEvent() {
	if f.submitButton.Button.Clicked() {
		if f.passwordEditor.IsValid() && f.passwordAgainEditor.IsValid() {
			passTxt := f.passwordEditor.Text()
			err := f.win.core.InitializeClient(passwordFromTxt(passTxt))
			if err == nil {
				f.win.Notify("Initialize app success")
				f.win.loggedIn = true
				f.callback()
			} else {
				f.win.Notify(fmt.Sprintf("Initialize app failed: %v", err))
			}
		}
	}
}
