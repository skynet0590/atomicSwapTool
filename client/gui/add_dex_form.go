package gui

import (
	"fmt"
	"gioui.org/layout"
	"github.com/skynet0590/atomicSwapTool/client/core"
	"github.com/skynet0590/atomicSwapTool/client/gui/dcrcomponent"
	"github.com/skynet0590/atomicSwapTool/client/gui/validate/validators"
	"github.com/skynet0590/atomicSwapTool/client/gui/values"
	"github.com/skynet0590/atomicSwapTool/dex"
	"io/ioutil"
)

type (
	dexServerForm struct {
		win            *Window
		dexAddress     *dcrcomponent.Editor
		tlsCertificate *dcrcomponent.Editor
		submitButton   dcrcomponent.Button

		passwordEditor *dcrcomponent.Editor
		confirmButton  dcrcomponent.Button

		callback    func()
		step        dexServerFormStep
		fee         uint64
		certContent []byte
	}
	dexServerFormStep uint8
)

const (
	dexServerInfo dexServerFormStep = 1 << iota
	dexServerConfirm
)

func newDEXServerForm(win *Window, callback func()) *dexServerForm {
	th := win.theme
	dexForm := &dexServerForm{
		win:            win,
		dexAddress:     th.Editor("DEX Address", validators.Required("DEX Address")),
		tlsCertificate: th.Editor("TLS Certificate", validators.Required("TLS Certificate")),
		submitButton:   th.Button("Add"),
		passwordEditor: th.EditorPassword("Password"),
		confirmButton:  th.Button("Confirm"),
		callback:       callback,
		step:           dexServerInfo,
	}

	return dexForm
}

func (f *dexServerForm) addDEXLayout(gtx C) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return f.dexAddress.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return f.tlsCertificate.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return f.submitButton.Layout(gtx)
		}),
	)
}

func (f *dexServerForm) confirmLayout(gtx C) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return f.passwordEditor.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return f.confirmButton.Layout(gtx)
		}),
	)
}

func (f *dexServerForm) Layout(gtx C) D {
	f.handlerEvent()
	title := "Add a DEX."
	if f.step == dexServerConfirm {
		title = fmt.Sprintf(
			`Enter your app password to confirm DEX registration. When you submit this form,
%d DCR will be spent from your Decred wallet to pay registration fees.`, f.fee)
	}
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
				return f.win.theme.H6(title).
					Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return f.win.theme.Line().Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			if f.step == dexServerInfo {
				return f.addDEXLayout(gtx)
			}
			if f.step == dexServerConfirm {
				return f.confirmLayout(gtx)
			}
			f.win.ChangePage(overview)
			return D{}
		}),
	)
}

func (f *dexServerForm) handlerEvent() {
	if f.submitButton.Button.Clicked() {
		certContent, err := f.getCertContent()
		if err != nil {
			f.win.Notify(fmt.Sprintf("Fail when open cert file: %v", err))
			return
		}
		fee, err := f.win.core.GetFee(f.dexAddress.Editor.Text(), certContent)
		if err != nil {
			f.win.Notify(err.Error())
		} else {
			f.step = dexServerConfirm
			f.fee = fee
		}
	}
	if f.confirmButton.Button.Clicked() {
		dcrID, _ := dex.BipSymbolID("dcr")
		wallet := f.win.core.WalletState(dcrID)
		if wallet == nil {
			f.win.Notify("No Decred wallet")
			return
		}

		_, err := f.win.core.Register(&core.RegisterForm{
			Addr:    f.dexAddress.Editor.Text(),
			Cert:    f.certContent,
			AppPass: passwordFromTxt(f.passwordEditor.Text()),
			Fee:     f.fee,
		})
		if err != nil {
			f.win.Notify("registration error: %v", err)
			return
		}
		f.callback()
	}
}

func (f *dexServerForm) getCertContent() (interface{}, error) {
	certPath := f.tlsCertificate.Text()
	content, err := ioutil.ReadFile(certPath)
	f.certContent = content
	return content, err
}
