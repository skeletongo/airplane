package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type boxMouseRender struct {
	fyne.Size
}

func (b *boxMouseRender) Destroy() {
}

func (b *boxMouseRender) Layout(size fyne.Size) {
}

func (b *boxMouseRender) MinSize() fyne.Size {
	return b.Size
}

func (b *boxMouseRender) Objects() []fyne.CanvasObject {
	return nil
}

func (b *boxMouseRender) Refresh() {
}

// BoxMouse 鼠标位置监听
type BoxMouse struct {
	widget.BaseWidget
}

func (b *BoxMouse) MouseDown(event *desktop.MouseEvent) {
	if event.Button != desktop.MouseButtonPrimary {
		return
	}
	Done(func() {
		// 判断点击位置是否在英雄机范围内
		hero := M.Airplanes[AirplaneHero][0]
		pos1 := hero.Position()
		pos2 := hero.Position().Add(hero.Size())

		if event.Position.X > pos1.X &&
			event.Position.X < pos2.X &&
			event.Position.Y > pos1.Y &&
			event.Position.Y < pos2.Y {
			M.HeroClick = true
			M.SubPos = fyne.NewPos(event.Position.X-pos1.X, event.Position.Y-pos1.Y)
		} else {
			M.HeroClick = false
		}
	})
}

func (b *BoxMouse) MouseUp(event *desktop.MouseEvent) {
	if event.Button != desktop.MouseButtonPrimary {
		return
	}
	Done(func() {
		M.HeroClick = false
	})
}

func (b *BoxMouse) CreateRenderer() fyne.WidgetRenderer {
	return &boxMouseRender{
		b.BaseWidget.Size(),
	}
}

func (b *BoxMouse) MouseIn(event *desktop.MouseEvent) {
}

func (b *BoxMouse) MouseMoved(event *desktop.MouseEvent) {
	if event.Button != desktop.MouseButtonPrimary {
		return
	}
	// 英雄机移动
	Done(func() {
		if M.Status != GameStatusGaming || !M.HeroClick {
			return
		}
		M.Airplanes[AirplaneHero][0].Move(event.Position.Subtract(M.SubPos))
	})
}

func (b *BoxMouse) MouseOut() {
}

func NewBoxMouse(size fyne.Size) *BoxMouse {
	ret := new(BoxMouse)
	ret.ExtendBaseWidget(ret)
	ret.Resize(size)
	return ret
}
