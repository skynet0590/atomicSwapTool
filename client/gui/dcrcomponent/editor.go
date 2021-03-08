// SPDX-License-Identifier: Unlicense OR MIT

package dcrcomponent

import (
	"github.com/skynet0590/atomicSwapTool/client/gui/validate"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"

	"github.com/planetdecred/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Clipboard interface {
	SetClipboard(e interface{})
	GetClipboard() string
}

type Editor struct {
	Clipboard Clipboard
	material.EditorStyle

	TitleLabel   Label
	ErrorLabel   Label
	LineColor    color.NRGBA
	DangerColor  color.NRGBA
	SurfaceColor color.NRGBA
	HintColor    color.NRGBA

	FlexWidth float32
	//IsVisible if true, displays the paste and clear button.
	IsVisible bool
	//IsTitleLabel if true makes the title label visible.
	IsTitleLabel bool
	//Bordered if true makes the adds a border around the editor.
	Bordered bool

	pasteBtnMaterial material.IconButtonStyle
	clearBtMaterial  material.IconButtonStyle

	Validators []validate.Validator

	M2 unit.Value
	M5 unit.Value
}

func (e Editor) Layout(gtx layout.Context) layout.Dimensions {
	e.handleEvents()

	if e.IsVisible {
		e.FlexWidth = 20
	}

	if e.Editor.Len() > 0 {
		e.TitleLabel.Text = e.Hint
	}

	c := color.NRGBA{R: 41, G: 112, B: 255, A: 255}
	if e.Editor.Focused() {
		e.TitleLabel.Text = e.Hint
		e.TitleLabel.Color = c
		e.LineColor = c
		e.Hint = ""
	}

	if e.ErrorLabel.Text != "" {
		e.LineColor, e.TitleLabel.Color = e.DangerColor, e.DangerColor
	}

	return layout.UniformInset(e.M2).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Top:    values.MarginPadding5,
					Bottom: values.MarginPadding5,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return e.editorLayout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									if e.ErrorLabel.Text != "" {
										inset := layout.Inset{
											Top:  e.M2,
											Left: e.M5,
										}
										return inset.Layout(gtx, func(gtx C) D {
											return e.ErrorLabel.Layout(gtx)
										})
									}
									return layout.Dimensions{}
								}),
							)
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							if e.IsTitleLabel {
								return layout.Inset{
									Top:  values.MarginPaddingMinus10,
									Left: values.MarginPadding10,
								}.Layout(gtx, func(gtx C) D {
									return Card{Color: e.SurfaceColor}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return e.TitleLabel.Layout(gtx)
									})
								})
							}
							return layout.Dimensions{}
						}),
					)
				})
			}),
		)
	})
}

func (e Editor) editorLayout(gtx C) D {
	if e.Bordered {
		border := widget.Border{Color: e.LineColor, CornerRadius: e.M5, Width: unit.Dp(1)}
		return border.Layout(gtx, func(gtx C) D {
			inset := layout.Inset{
				Top:    e.M2,
				Bottom: e.M2,
				Left:   values.MarginPadding10,
				Right:  e.M5,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return e.editor(gtx)
			})
		})
	}

	return e.editor(gtx)
}

func (e Editor) editor(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{
						Top:    e.M5,
						Bottom: e.M5,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return e.EditorStyle.Layout(gtx)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			if e.IsVisible {
				inset := layout.Inset{
					Top:  e.M2,
					Left: e.M5,
				}
				return inset.Layout(gtx, func(gtx C) D {
					if e.Editor.Text() == "" {
						return e.pasteBtnMaterial.Layout(gtx)
					}
					return e.clearBtMaterial.Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
	)
}

func (e *Editor) validate() {
	for _, v := range e.Validators {
		valid, errTxt := v.Validate(e.Editor.Text())
		if !valid {
			e.ErrorLabel.Text = errTxt
			return
		}
	}
	e.ClearError()
}

func (e *Editor) handleEvents() {
	e.validate()
	if e.pasteBtnMaterial.Button.Clicked() {
		e.Editor.Focus()

		go func() {
			text := e.Clipboard.GetClipboard()
			e.Editor.SetText(text)
			//e.Editor.Move(e.Editor.Len())
		}()
		go func() {
			e.Clipboard.SetClipboard(ReadClipboard{})
		}()
	}

	for e.clearBtMaterial.Button.Clicked() {
		e.Editor.SetText("")
	}

	if e.ErrorLabel.Text != "" {
		e.LineColor = e.DangerColor
	} else {
		// e.LineColor = e.HintColor
	}
}

func (e *Editor) SetError(text string) {
	e.ErrorLabel.Text = text
}

func (e *Editor) ClearError() {
	e.ErrorLabel.Text = ""
}

func (e *Editor) IsDirty() bool {
	return e.ErrorLabel.Text == ""
}

func (t *Theme) Editor(hint string, validators ...validate.Validator) Editor {
	editor := new(widget.Editor)
	editor.SingleLine = true
	errorLabel := t.Caption("")
	errorLabel.Color = t.Color.Danger

	m := material.Editor(t.Theme, editor, hint)
	m.TextSize = t.TextSize
	m.Color = t.Color.Text
	m.Hint = hint
	m.HintColor = t.Color.Hint

	var m0 = unit.Dp(0)
	var m25 = unit.Dp(25)

	e := Editor{
		Clipboard:    t.Clipboard,
		EditorStyle:  m,
		TitleLabel:   t.Body2(""),
		FlexWidth:    0,
		IsTitleLabel: true,
		Bordered:     true,
		LineColor:    t.Color.Hint,
		HintColor:    t.Color.Hint,
		DangerColor:  t.Color.Danger,
		SurfaceColor: t.Color.Surface,
		ErrorLabel:   errorLabel,

		M2: unit.Dp(2),
		M5: unit.Dp(5),

		pasteBtnMaterial: material.IconButtonStyle{
			Icon:       mustIcon(widget.NewIcon(icons.ContentContentPaste)),
			Size:       m25,
			Background: color.NRGBA{},
			Color:      t.Color.Text,
			Inset:      layout.UniformInset(m0),
			Button:     new(widget.Clickable),
		},

		clearBtMaterial: material.IconButtonStyle{
			Icon:       mustIcon(widget.NewIcon(icons.ContentClear)),
			Size:       m25,
			Background: color.NRGBA{},
			Color:      t.Color.Text,
			Inset:      layout.UniformInset(m0),
			Button:     new(widget.Clickable),
		},
		Validators: validators,
	}

	return e
}
