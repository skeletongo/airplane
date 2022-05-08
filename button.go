package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type PauseButton struct {
	widget.Icon
	flag1, flag2 bool
	mouseDown    func()
	mouseUp      func()
}

func (s *PauseButton) MouseDown(event *desktop.MouseEvent) {
	if s.flag1 {
		s.SetResource(resourceResumepressedPng)
	} else {
		s.SetResource(resourcePausepressedPng)
	}
	s.flag1 = !s.flag1
	s.mouseDown()
}

func (s *PauseButton) MouseUp(event *desktop.MouseEvent) {
	if s.flag2 {
		s.SetResource(resourcePausenorPng)
	} else {
		s.SetResource(resourceResumenorPng)
	}
	s.flag2 = !s.flag2
	s.mouseUp()
}

func (s *PauseButton) ToDefault() {
	s.flag1 = false
	s.flag2 = false
	s.Icon.SetResource(resourcePausenorPng)
	s.Resize(fyne.NewSize(65, 45))
}

func NewPauseButton(mouseDown func(), mouseUp func()) *PauseButton {
	ret := new(PauseButton)
	ret.ExtendBaseWidget(ret)
	ret.ToDefault()
	ret.mouseDown = mouseDown
	ret.mouseUp = mouseUp
	return ret
}

type BombButton struct {
	widget.Icon

	onTapped func()
}

func (b *BombButton) Tapped(event *fyne.PointEvent) {
	b.onTapped()
}

func NewBombButton(tapped func()) *BombButton {
	ret := new(BombButton)
	ret.ExtendBaseWidget(ret)
	ret.SetResource(resourceBombPng)
	ret.Resize(fyne.NewSize(63, 57))
	ret.Move(fyne.NewPos(MasterWidth-63-10, MasterHeight-57-10))
	ret.onTapped = tapped
	return ret
}

type ItemButton struct {
	widget.Icon

	onTapped func()
}

func (i *ItemButton) Tapped(event *fyne.PointEvent) {
	i.onTapped()
}

func (i *ItemButton) MinSize() fyne.Size {
	return fyne.NewSize(300, 41)
}

func NewItemButton(res fyne.Resource, tapped func()) *ItemButton {
	ret := new(ItemButton)
	ret.ExtendBaseWidget(ret)
	ret.SetResource(res)
	ret.onTapped = tapped
	return ret
}
