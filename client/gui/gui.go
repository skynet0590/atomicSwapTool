// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package gui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/skynet0590/atomicSwapTool/client/core"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
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
		Color       Color
		TextSize              unit.Value

		Clipboard     chan string
		ReadClipboard chan interface{}
	}
	Page interface {
		Layout(gtx C) D
		Title() string
	}
	PageTitle string
	Color  struct {
		Primary    color.NRGBA
		Secondary  color.NRGBA
		Text       color.NRGBA
		Hint       color.NRGBA
		Overlay    color.NRGBA
		InvText    color.NRGBA
		Success    color.NRGBA
		Danger     color.NRGBA
		Background color.NRGBA
		Surface    color.NRGBA
		Gray       color.NRGBA
		Black      color.NRGBA
		Orange     color.NRGBA
		LightGray  color.NRGBA
	}
)

var (
	overview PageTitle = "overview"
	register PageTitle = "register"
)

// NewUI creates a new UI
func NewWindow(coreClient *core.Core) *Window {
	w := &Window{
		currentPage: register,
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

	w.initProperties()

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
		register: newRegisterPage(w),
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
					user := w.core.User()

					if user.Initialized {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Max.X /= 3
								return w.nav.Layout(gtx, w.theme, w.navAnim)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return page.Layout(gtx)
							}),
						)
					}else{
						page = w.pages[register]
						return layout.Flex{}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return page.Layout(gtx)
							}),
						)
					}
				})
				bar := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// w.bar.Title = page.Title()
					return w.bar.Layout(gtx, w.theme)
				})
				layout.Stack{
					Alignment: layout.N,
				}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return fill(gtx, w.Color.Background)
					}),
					layout.Stacked(func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx, bar, content)
					}),
				)
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

func (w *Window) initProperties() {
	c := Color{
		Primary:    rgb(0x2970ff), // key blue
		Secondary:  rgb(0x091440), // dark blue
		Text:       rgb(0x000000),
		Hint:       rgb(0xbbbbbb),
		Overlay:    rgb(0x000000),
		InvText:    rgb(0xffffff),
		Success:    rgb(0x41bf53),
		Danger:     rgb(0xff0000),
		Background: argb(0x22444444),
		Surface:    rgb(0xffffff),
		Gray:       rgb(0x596D81),
		Black:      rgb(0x000000),
		Orange:     rgb(0xed6d47),
		LightGray:  rgb(0xc4cbd2),
	}
	w.Color = c
	w.TextSize = unit.Sp(16)
	w.Clipboard = make(chan string)
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func fillMax(gtx layout.Context, col color.NRGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Max.X, Y: cs.Max.Y}
	st := op.Save(gtx.Ops)
	track := image.Rectangle{
		Max: image.Point{X: d.X, Y: d.Y},
	}
	clip.Rect(track).Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	st.Load()
}

func fill(gtx layout.Context, col color.NRGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := image.Point{X: cs.Min.X, Y: cs.Min.Y}
	st := op.Save(gtx.Ops)
	track := image.Rectangle{
		Max: d,
	}
	clip.Rect(track).Add(gtx.Ops)
	paint.Fill(gtx.Ops, col)
	st.Load()

	return layout.Dimensions{Size: d}
}

func (w *Window) GetClipboard() string {
	txt := <-w.Clipboard
	return txt
}

func (w *Window) SetClipboard(e interface{}) {
	w.ReadClipboard <- e
}
