package gui

import (
	"gioui.org/layout"
	"github.com/skynet0590/atomicSwapTool/client/gui/dcrcomponent"
	"github.com/skynet0590/atomicSwapTool/client/gui/validate/validators"
	"github.com/skynet0590/atomicSwapTool/client/gui/values"
)

type (
	dcrWalletSettingForm struct {
		win                            *Window
		accountNameEditor              *dcrcomponent.Editor
		rpcUserNameEditor              *dcrcomponent.Editor
		rpcPasswordEditor              *dcrcomponent.Editor
		rpcAddressEditor               *dcrcomponent.Editor
		tlsCertificateEditor           *dcrcomponent.Editor
		fallbackFeeRateEditor          *dcrcomponent.Editor
		redeemConfirmationTargetEditor *dcrcomponent.Editor
		walletPasswordEditor           *dcrcomponent.Editor
		appPasswordEditor              *dcrcomponent.Editor
		submitButton                   dcrcomponent.Button
	}
)

func newDcrWalletSetting(w *Window) *dcrWalletSettingForm {
	th := w.theme
	return &dcrWalletSettingForm{
		win:                            w,
		accountNameEditor:              th.Editor("dcrwallet account name", validators.Required("Account name")),
		rpcUserNameEditor:              th.Editor("RPC Username",validators.Required("RPC Username")),
		rpcPasswordEditor:              th.EditorPassword("RPC Password", validators.Required("RPC Password")),
		rpcAddressEditor:               th.Editor("RPC Address"),
		tlsCertificateEditor:           th.Editor("RPC Certificate"),
		fallbackFeeRateEditor:          th.Editor("Fallback fee rate"),
		redeemConfirmationTargetEditor: th.Editor("Redeem confirmation target"),
		walletPasswordEditor:           th.EditorPassword("Wallet Password"),
		appPasswordEditor:              th.EditorPassword("App Password"),

	}
}

func (f *dcrWalletSettingForm) Layout(gtx C) D {
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
				return f.win.theme.H6("Your Decred wallet is required to pay registration fees.").
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
						return f.accountNameEditor.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return f.rpcUserNameEditor.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return f.rpcPasswordEditor.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return f.rpcAddressEditor.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return f.tlsCertificateEditor.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return f.fallbackFeeRateEditor.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return f.redeemConfirmationTargetEditor.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return f.walletPasswordEditor.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return f.submitButton.Layout(gtx)
					}),
				)
			})
		}),
	)
}
