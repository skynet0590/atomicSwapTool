// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package gui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/skynet0590/atomicSwapTool/client/core"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"time"
)

// defaultMargin is a margin applied in multiple places to give
// widgets room to breathe.
var defaultMargin = unit.Dp(10)

type (
	C = layout.Context
	D = layout.Dimensions
	// Window holds all of the application state.
	Window struct {
		core        *core.Core
		currentPage PageTitle
		theme       *material.Theme
		window      *app.Window
		ops         *op.Ops
		modal       *component.ModalLayer
		pages       map[PageTitle]Page
		bar         *component.AppBar
		nav         *component.NavDrawer
		navAnim     *component.VisibilityAnimation
	}
	Page interface {
		Layout(gtx C) D
		Title() string
	}
	PageTitle string
)

var (
	overview PageTitle = "overview"
)

// NewUI creates a new UI

func NewWindow(coreClient *core.Core) *Window {
	w := &Window{
		currentPage: overview,
		core:        coreClient,
	}
	w.theme = material.NewTheme(gofont.Collection())
	w.window = app.NewWindow(app.Title("Atomic Swap Client"))
	w.ops = &op.Ops{}
	w.modal = component.NewModal()

	// Setup app bar
	w.bar = component.NewAppBar(w.modal)
	w.bar.Title = "Atomic Swap Client"
	icon, _ := widget.NewIcon(icons.ActionHome)
	w.bar.NavigationIcon = icon

	// Setup left nav
	nav := component.NewNav("Navigation Drawer", "")

	w.nav = &nav
	navAnim := component.VisibilityAnimation{
		Duration: time.Millisecond * 100,
		State:    component.Invisible,
	}
	w.navAnim = &navAnim
	w.navAnim.State = component.Appearing
	w.pages = map[PageTitle]Page{
		overview: newOverviewPage(w),
	}
	return w
}

// Run handles window events and renders the application.
func (w *Window) Loop() error {
	win := w.window
	for {
		select {
		case e := <-win.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				return evt.Err
			case system.FrameEvent:
				gtx := layout.NewContext(w.ops, evt)
				page := w.pages[w.currentPage]
				content := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max.X /= 3
							return w.nav.Layout(gtx, w.theme, w.navAnim)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return page.Layout(gtx)
						}),
					)
				})
				bar := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// w.bar.Title = page.Title()
					return w.bar.Layout(gtx, w.theme)
				})
				layout.Flex{Axis: layout.Vertical}.Layout(gtx, bar, content)
				evt.Frame(gtx.Ops)
			}
		}
	}
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
