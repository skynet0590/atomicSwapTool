// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package gui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// defaultMargin is a margin applied in multiple places to give
// widgets room to breathe.
var defaultMargin = unit.Dp(10)

// Window holds all of the application state.
type Window struct {
	Theme  *material.Theme
	window *app.Window
	ops    *op.Ops
}

// NewUI creates a new UI

func NewWindow() *Window {
	w := &Window{}
	w.Theme = material.NewTheme(gofont.Collection())
	w.window = app.NewWindow(app.Title("Atomic Swap Client"))
	w.ops = &op.Ops{}
	return w
}

// Run handles window events and renders the application.
func (w *Window) Loop() error {
	for {
		select {
		case e := <-w.window.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				return evt.Err
			case system.FrameEvent:
			}
		}
	}
	var ops op.Ops

	// listen for events happening on the window.
	for e := range w.window.Events() {
		// detect the type of the event.
		switch e := e.(type) {
		// this is sent when the application should re-render.
		case system.FrameEvent:
			// gtx is used to pass around rendering and event information.
			gtx := layout.NewContext(&ops, e)
			// render and handle UI.
			w.Layout(gtx)
			// render and handle the operations from the UI.
			e.Frame(gtx.Ops)

		// handle a global key press.
		case key.Event:
			switch e.Name {
			// when we click escape, let's close the window.
			case key.NameEscape:
				return nil
			}

		// this is sent when the application is closed.
		case system.DestroyEvent:
			return e.Err
		}
	}

	return nil
}

// Layout displays the main program layout.
func (w *Window) Layout(gtx layout.Context) layout.Dimensions {
	// inset is used to add padding around the window border.
	inset := layout.UniformInset(defaultMargin)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{}
		// return ui.Converter.Layout(ui.Theme, gtx)
	})
}
