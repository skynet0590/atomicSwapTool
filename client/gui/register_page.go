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
		registerForm *registerForm
		dcrWalletForm *dcrWalletSettingForm
		step         registerStep
	}
	registerForm struct {
		win                 *Window
		parents             *registerPage
		passwordEditor      *dcrcomponent.Editor
		passwordAgainEditor *dcrcomponent.Editor
		submitButton        dcrcomponent.Button
	}
	registerStep int
)

const (
	registerSetPassword registerStep = 0
	registerDCRWallet
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
			return p.registerForm.Layout(gtx)
		})
	})
}

func newRegisterForm(w *Window) *registerForm {
	const passTxt = "Password"
	passwordEditor := w.theme.EditorPassword(passTxt, validators.Required(passTxt))
	passwordAgainEditor := w.theme.EditorPassword("Password Again", validators.MatchedInput(passTxt, passwordEditor.Editor))
	return &registerForm{
		win:                 w,
		passwordEditor:      passwordEditor,
		passwordAgainEditor: passwordAgainEditor,
		submitButton:        w.theme.Button("Submit"),
	}
}

func newRegisterPage(w *Window) *registerPage {
	page := registerPage{
		win:          w,
		registerForm: newRegisterForm(w),
		dcrWalletForm: newDcrWalletSetting(w),
		step: registerSetPassword,
	}
	page.registerForm.parents = &page
	return &page
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
			err := f.win.core.InitializeClient([]byte(passTxt))
			if err == nil {
				f.win.Notify("Initialize app success")
				f.parents.step = registerDCRWallet
			}else{
				f.win.Notify(fmt.Sprintf("Initialize app faild: %v", err))
			}
		}
	}
}
