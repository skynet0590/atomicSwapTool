// SPDX-License-Identifier: Unlicense OR MIT

package dcrcomponent

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"gioui.org/widget"
)

type Button struct {
	material.ButtonStyle
}

type IconButton struct {
	material.IconButtonStyle
}

func Clickable(gtx layout.Context, button *widget.Clickable, w layout.Widget) layout.Dimensions {
	return material.Clickable(gtx, button, w)
}

func (b Button) Layout(gtx layout.Context) layout.Dimensions {
	return b.ButtonStyle.Layout(gtx)
}

func (b IconButton) Layout(gtx layout.Context) layout.Dimensions {
	return b.IconButtonStyle.Layout(gtx)
}

type TextAndIconButton struct {
	theme           *material.Theme
	Button          *widget.Clickable
	icon            *widget.Icon
	text            string
	Color           color.NRGBA
	BackgroundColor color.NRGBA
}

func NewTextAndIconButton(theme *material.Theme, btn *widget.Clickable, icon *widget.Icon, txt string) TextAndIconButton {
	return TextAndIconButton {
		theme:           theme,
		Button:          btn,
		icon:            icon,
		text:            txt,
	}
}

func (b TextAndIconButton) Layout(gtx layout.Context) layout.Dimensions {
	btnLayout := material.ButtonLayout(b.theme, b.Button)
	btnLayout.Background = b.BackgroundColor
	b.icon.Color = b.Color

	return btnLayout.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx C) D {
			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			textIconSpacer := unit.Dp(5)

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					var d D
					size := gtx.Px(unit.Dp(46)) - 2*gtx.Px(unit.Dp(16))
					b.icon.Layout(gtx, unit.Px(float32(size)))
					d = layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					l := material.Body1(b.theme, b.text)
					l.Color = b.Color
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layLabel, layIcon)
		})
	})
}
