package dcrcomponent

import (
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
)

func (t *Theme) Card() Card {
	return Card{
		Color: t.Color.Surface,
		Radius: CornerRadius{
			NE: DefaultRadius,
			SE: DefaultRadius,
			NW: DefaultRadius,
			SW: DefaultRadius,
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
func (t *Theme) Line() *Line {
	col := t.Color.Primary
	col.A = 150

	return &Line{
		Height: 1,
		Color:  col,
	}
}

func (t *Theme) Button(txt string) Button {
	return Button{material.Button(t.Theme, new(widget.Clickable), txt)}
}

func (t *Theme) IconButton(icon *widget.Icon) IconButton {
	return IconButton{material.IconButton(t.Theme, new(widget.Clickable), icon)}
}

func (t *Theme) PlainIconButton(button *widget.Clickable, icon *widget.Icon) IconButton {
	btn := IconButton{material.IconButton(t.Theme, button, icon)}
	btn.Background = color.NRGBA{}
	return btn
}

func (t *Theme) TextAndIconButton(text string, icon *widget.Icon) TextAndIconButton {
	btn := NewTextAndIconButton(t.Theme, new(widget.Clickable), icon, text)
	btn.Color = t.Color.Surface
	btn.BackgroundColor = t.Color.Primary
	return btn
}
