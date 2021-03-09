// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package gui

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/x/component"
	"gioui.org/x/notify"
	"github.com/skynet0590/atomicSwapTool/client/core"
	"github.com/skynet0590/atomicSwapTool/client/gui/dcrcomponent"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
	"log"
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
		theme       *dcrcomponent.Theme
		window      *app.Window
		ops         *op.Ops
		modal       *component.ModalLayer
		pages       map[PageTitle]Page
		bar         *component.AppBar
		nav         *component.NavDrawer
		navAnim     *component.VisibilityAnimation
		TextSize    unit.Value

		notify        notify.Manager
		Clipboard     chan string
		ReadClipboard chan interface{}
	}
	Page interface {
		Layout(gtx C) D
		Title() string
	}
	PageTitle string
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
	w.theme = dcrcomponent.NewTheme(gofont.Collection(), dcrcomponent.DefaultColor, w)
	w.window = app.NewWindow(app.Title("Atomic Swap Client"))
	w.ops = &op.Ops{}
	w.modal = component.NewModal()

	// Setup app bar
	w.bar = component.NewAppBar(w.modal)
	w.bar.Title = "Atomic Swap Client"
	icon, _ := widget.NewIcon(icons.ActionHome)
	w.bar.NavigationIcon = icon
	notif,err := notify.NewManager()
	if err != nil {
		fmt.Println("Init notify failed: ", err)
	}
	w.notify = notif
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
								return w.nav.Layout(gtx, w.theme.Theme, w.navAnim)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return page.Layout(gtx)
							}),
						)
					} else {
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
					return w.bar.Layout(gtx, w.theme.Theme)
				})
				layout.Stack{
					Alignment: layout.N,
				}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return dcrcomponent.Fill(gtx, w.theme.Color.Background)
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
	w.TextSize = unit.Sp(16)
	w.Clipboard = make(chan string)
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func (w *Window) GetClipboard() string {
	txt := <-w.Clipboard
	return txt
}

func (w *Window) SetClipboard(e interface{}) {
	w.ReadClipboard <- e
}

func (w *Window) Notify(txt string) {
	notif, e := w.notify.CreateNotification("Atomic Swap", txt)
	if e != nil {
		log.Printf("notification send failed: %v", e)
	}
	go func() {
		time.Sleep(time.Second * 10)
		if e = notif.Cancel(); e != nil {
			log.Printf("failed cancelling: %v", e)
		}
	}()
}

func (w *Window) ChangePage(page PageTitle) {
	w.currentPage = page
}
