package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

const (
	PropBomb = iota
	PropBullet
	PropMax
)

const (
	PropStatusDefault = iota
	PropStatusFlying
)

type Prop struct {
	fyne.CanvasObject
	Tp     int
	Status int

	speed float32 // 移动速度,单位 px/s
}

func (p *Prop) ToDefault() {
	p.Status = PropStatusDefault
	p.Move(fyne.NewPos(0, -p.Size().Height))
}

func NewProp(tp int, status int) *Prop {
	ret := &Prop{
		Tp:     tp,
		Status: status,
	}
	return ret
}

func AddProp(c *fyne.Container, prop *Prop) {
	var img *canvas.Image
	if prop.Tp == PropBomb {
		img = getImage(resourceBombsupplyPng, 60, 107)
	} else {
		img = getImage(resourceBulletsupplyPng, 58, 88)
	}
	prop.CanvasObject = img
	c.Add(img)

	prop.ToDefault()
}
