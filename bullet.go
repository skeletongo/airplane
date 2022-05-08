package main

import (
	"fyne.io/fyne/v2"
)

const (
	BulletHero  = iota
	BulletEnemy //todo 敌机子弹
	BulletMax
)

const (
	BulletStatusUnused = iota
	BulletStatusFlying
)

const (
	BulletCreateHero = iota
	BulletCreateEnemy
)

type Bullet struct {
	fyne.CanvasObject
	Tp     int // 子弹是谁发射的
	Status int // 子弹状态,飞行中或未使用
}

func (b *Bullet) ToDefault() {
	b.Status = BulletStatusUnused
	b.Move(fyne.NewPos(0, -b.Size().Height))
}

// NewBullet 创建子弹
// tp 子弹
// status 状态
func NewBullet(tp int, status int) *Bullet {
	ret := &Bullet{
		Tp:     tp,
		Status: status,
	}
	return ret
}

func AddBullet(c *fyne.Container, bullet *Bullet) {
	img := getImage(resourceBullet2Png, 5, 11)
	bullet.CanvasObject = img
	c.Add(img)

	bullet.ToDefault()
}
