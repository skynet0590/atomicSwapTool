package gui

import (
	"gioui.org/layout"
	"github.com/skynet0590/atomicSwapTool/client/gui/dcrcomponent"
	"github.com/skynet0590/atomicSwapTool/client/gui/values"
)

type confirmForm struct {
	win *Window
	title *string
	passwordEditor      *dcrcomponent.Editor
	submitButton        dcrcomponent.Button
}

func newConfirmForm(w *Window, title *string, confirmButtonTxt string) *confirmForm {
	if confirmButtonTxt == "" {
		confirmButtonTxt = "Submit"
	}
	th := w.theme
	form := &confirmForm{
		win: w,
		title: title,
		passwordEditor: th.EditorPassword("Password"),
		submitButton: th.Button(confirmButtonTxt),
	}
	return form
}

func (f *confirmForm) Layout(gtx C) D {
	f.handleEvents()
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
				return f.win.theme.H6(*f.title).
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
						return f.submitButton.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (f *confirmForm) handleEvents() {

}
