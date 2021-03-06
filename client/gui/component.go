package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/skynet0590/atomicSwapTool/client/gui/dcrcomponent"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

func (w *Window) Card() dcrcomponent.Card {
	return dcrcomponent.Card{
		Color: w.Color.Surface,
		Radius: dcrcomponent.CornerRadius{
			NE: dcrcomponent.DefaultRadius,
			SE: dcrcomponent.DefaultRadius,
			NW: dcrcomponent.DefaultRadius,
			SW: dcrcomponent.DefaultRadius,
		},
	}
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

// Line returns a line widget instance
func (w *Window) Line() *dcrcomponent.Line {
	col := w.Color.Primary
	col.A = 150

	return &dcrcomponent.Line{
		Height: 1,
		Color:  col,
	}
}

func (w *Window) H1(txt string) material.LabelStyle {
	return material.H1(w.theme, txt)
}

func (w *Window) H2(txt string) material.LabelStyle {
	return material.H2(w.theme, txt)
}

func (w *Window) H3(txt string) material.LabelStyle {
	return material.H3(w.theme, txt)
}

func (w *Window) H4(txt string) material.LabelStyle {
	return material.H4(w.theme, txt)
}

func (w *Window) H5(txt string) material.LabelStyle {
	return material.H5(w.theme, txt)
}

func (w *Window) H6(txt string) material.LabelStyle {
	return material.H6(w.theme, txt)
}

func (w *Window) Body1(txt string) material.LabelStyle {
	return material.Body1(w.theme, txt)
}

func (w *Window) Body2(txt string) material.LabelStyle {
	return material.Body2(w.theme, txt)
}

func (w *Window) Caption(txt string) material.LabelStyle {
	return material.Caption(w.theme, txt)
}

func (w *Window) Editor(hint string) dcrcomponent.Editor {
	editor := new(widget.Editor)
	editor.SingleLine = true
	errorLabel := w.Caption("")
	errorLabel.Color = w.Color.Danger

	m := material.Editor(w.theme, editor, hint)
	m.TextSize = w.TextSize
	m.Color = w.Color.Text
	m.Hint = hint
	m.HintColor = w.Color.Hint

	var m0 = unit.Dp(0)
	var m25 = unit.Dp(25)

	e := dcrcomponent.Editor{
		Clipboard: w,
		EditorStyle:       m,
		TitleLabel:        w.Body2(""),
		FlexWidth:         0,
		IsTitleLabel:      true,
		Bordered:          true,
		LineColor:         w.Color.Hint,
		ErrorLabel:        errorLabel,
		RequiredErrorText: "Field is required",

		M2: unit.Dp(2),
		M5: unit.Dp(5),

		PasteBtnMaterial: material.IconButtonStyle{
			Icon:       mustIcon(widget.NewIcon(icons.ContentContentPaste)),
			Size:       m25,
			Background: color.NRGBA{},
			Color:      w.Color.Text,
			Inset:      layout.UniformInset(m0),
			Button:     new(widget.Clickable),
		},

		ClearBtMaterial: material.IconButtonStyle{
			Icon:       mustIcon(widget.NewIcon(icons.ContentClear)),
			Size:       m25,
			Background: color.NRGBA{},
			Color:      w.Color.Text,
			Inset:      layout.UniformInset(m0),
			Button:     new(widget.Clickable),
		},
	}

	return e
}
