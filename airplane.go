package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"time"
)

const (
	AirplaneEnemy1 = iota // 小型敌机
	AirplaneEnemy2        // 中型敌机
	AirplaneEnemy3        // 大型敌机
	AirplaneHero          // 英雄机
	AirplaneMax
)

const (
	AirplaneStatusDefault = iota // 没有使用的状态
	AirplaneStatusFlying         // 飞行中
	AirplaneStatusHit            // 被子弹击中
	AirplaneStatusDown           // 爆炸中
	AirplaneStatusMax
)

var AirplaneImages [][][]*canvas.Image

func getImage(r *fyne.StaticResource, w, h float32) *canvas.Image {
	ret := canvas.NewImageFromResource(r)
	ret.SetMinSize(fyne.NewSize(w, h))
	ret.Resize(ret.MinSize())
	return ret
}

func init() {
	// 飞机不同状态的动画图片数据表
	AirplaneImages = make([][][]*canvas.Image, AirplaneMax)
	AirplaneImages[AirplaneEnemy1] = make([][]*canvas.Image, AirplaneStatusMax)
	AirplaneImages[AirplaneEnemy2] = make([][]*canvas.Image, AirplaneStatusMax)
	AirplaneImages[AirplaneEnemy3] = make([][]*canvas.Image, AirplaneStatusMax)
	AirplaneImages[AirplaneHero] = make([][]*canvas.Image, AirplaneStatusMax)

	AirplaneImages[AirplaneEnemy1][AirplaneStatusDefault] = []*canvas.Image{
		getImage(resourceEnemy1Png, 57, 43),
	}
	AirplaneImages[AirplaneEnemy1][AirplaneStatusFlying] = []*canvas.Image{
		getImage(resourceEnemy1Png, 57, 43),
	}
	AirplaneImages[AirplaneEnemy1][AirplaneStatusHit] = []*canvas.Image{
		getImage(resourceEnemy1Png, 57, 43),
	}
	AirplaneImages[AirplaneEnemy1][AirplaneStatusDown] = []*canvas.Image{
		getImage(resourceEnemy1down1Png, 57, 51),
		getImage(resourceEnemy1down2Png, 57, 51),
		getImage(resourceEnemy1down3Png, 57, 51),
		getImage(resourceEnemy1down4Png, 57, 51),
	}

	AirplaneImages[AirplaneEnemy2][AirplaneStatusDefault] = []*canvas.Image{
		getImage(resourceEnemy2Png, 69, 99),
	}
	AirplaneImages[AirplaneEnemy2][AirplaneStatusFlying] = []*canvas.Image{
		getImage(resourceEnemy2Png, 69, 99),
	}
	AirplaneImages[AirplaneEnemy2][AirplaneStatusHit] = []*canvas.Image{
		getImage(resourceEnemy2hitPng, 69, 99),
	}
	AirplaneImages[AirplaneEnemy2][AirplaneStatusDown] = []*canvas.Image{
		getImage(resourceEnemy2down1Png, 69, 95),
		getImage(resourceEnemy2down2Png, 69, 95),
		getImage(resourceEnemy2down3Png, 69, 95),
		getImage(resourceEnemy2down4Png, 69, 95),
	}

	AirplaneImages[AirplaneEnemy3][AirplaneStatusDefault] = []*canvas.Image{
		getImage(resourceEnemy3n1Png, 169, 258),
	}
	AirplaneImages[AirplaneEnemy3][AirplaneStatusFlying] = []*canvas.Image{
		getImage(resourceEnemy3n1Png, 169, 258),
		getImage(resourceEnemy3n2Png, 169, 258),
	}
	AirplaneImages[AirplaneEnemy3][AirplaneStatusHit] = []*canvas.Image{
		getImage(resourceEnemy3hitPng, 169, 258),
	}
	AirplaneImages[AirplaneEnemy3][AirplaneStatusDown] = []*canvas.Image{
		getImage(resourceEnemy3down1Png, 165, 261),
		getImage(resourceEnemy3down2Png, 165, 261),
		getImage(resourceEnemy3down3Png, 165, 261),
		getImage(resourceEnemy3down4Png, 165, 261),
		getImage(resourceEnemy3down5Png, 165, 261),
		getImage(resourceEnemy3down6Png, 165, 261),
	}

	AirplaneImages[AirplaneHero][AirplaneStatusDefault] = []*canvas.Image{
		getImage(resourceMe1Png, 102, 126),
	}
	AirplaneImages[AirplaneHero][AirplaneStatusFlying] = []*canvas.Image{
		getImage(resourceMe1Png, 102, 126),
		getImage(resourceMe2Png, 102, 126),
	}
	AirplaneImages[AirplaneHero][AirplaneStatusHit] = []*canvas.Image{
		getImage(resourceMe1Png, 102, 126),
		getImage(resourceMe2Png, 102, 126),
	}
	AirplaneImages[AirplaneHero][AirplaneStatusDown] = []*canvas.Image{
		getImage(resourceMedestroy1Png, 102, 126),
		getImage(resourceMedestroy2Png, 102, 126),
		getImage(resourceMedestroy3Png, 102, 126),
		getImage(resourceMedestroy4Png, 102, 126),
	}
}

// Airplane 飞机对象
type Airplane struct {
	fyne.CanvasObject
	Tp     int // 飞机类型
	Status int // 飞机状态
	HP     int // 血量

	anim           [AirplaneStatusMax]*fyne.Animation
	animId         int
	speed          float32   // 移动速度,单位 px/s
	lastBulletTime time.Time // 最后一个子弹的发射时间
}

func (a *Airplane) ToAnim(i int) {
	if i < 0 || i >= AirplaneStatusMax {
		return
	}
	if a.anim[a.animId] != nil {
		a.anim[a.animId].Stop()
	}
	if a.anim[i] != nil {
		a.anim[i].Start()
	} else {
		box := a.CanvasObject.(*fyne.Container)
		img := AirplaneImages[a.Tp][a.Status][0]
		box.Objects = []fyne.CanvasObject{img}
		box.Resize(img.Size())
		box.Refresh()
	}
	a.animId = i
}

func (a *Airplane) ToHit() {
	// 敌机,击中动画，扣血量
	a.ToAnim(AirplaneStatusHit)
	a.HP--
	if a.HP <= 0 {
		a.Status = AirplaneStatusDown
		a.ToAnim(AirplaneStatusDown)
	}
}

// ToDefault 恢复初始状态
func (a *Airplane) ToDefault() {
	a.Status = AirplaneStatusDefault
	switch a.Tp {
	case AirplaneHero:
		a.ToAnim(AirplaneStatusFlying)
		a.Move(fyne.NewPos((MasterWidth-a.Size().Width)/2, MasterHeight-200))
		a.HP = 1
	case AirplaneEnemy1:
		a.ToAnim(AirplaneStatusDefault)
		a.Move(fyne.NewPos(0, -a.Size().Height))
		a.HP = 1
	case AirplaneEnemy2:
		a.ToAnim(AirplaneStatusDefault)
		a.Move(fyne.NewPos(0, -a.Size().Height))
		a.HP = 3
	case AirplaneEnemy3:
		a.ToAnim(AirplaneStatusDefault)
		a.Move(fyne.NewPos(0, -a.Size().Height))
		a.HP = 5
	}
}

func NewAirplane(tp int, status int) *Airplane {
	ret := &Airplane{
		Tp:     tp,
		Status: status,
		anim:   [AirplaneStatusMax]*fyne.Animation{},
	}
	return ret
}

// AddAirplane 绘制飞机
func AddAirplane(c *fyne.Container, airplane *Airplane) {
	airplane.CanvasObject = container.NewWithoutLayout()
	airplane.ToDefault()
	// 动画
	switch airplane.Tp {
	case AirplaneEnemy1:
		CreateAirplaneAnim(airplane, time.Second, AirplaneStatusDown, 0)

	case AirplaneEnemy2:
		CreateAirplaneAnim(airplane, canvas.DurationShort, AirplaneStatusHit, 0)
		CreateAirplaneAnim(airplane, time.Second, AirplaneStatusDown, 0)

	case AirplaneEnemy3:
		CreateAirplaneAnim(airplane, canvas.DurationShort, AirplaneStatusFlying, fyne.AnimationRepeatForever)
		CreateAirplaneAnim(airplane, canvas.DurationShort, AirplaneStatusHit, 0)
		CreateAirplaneAnim(airplane, time.Second, AirplaneStatusDown, 0)

	case AirplaneHero:
		CreateAirplaneAnim(airplane, canvas.DurationShort, AirplaneStatusFlying, fyne.AnimationRepeatForever)
		CreateAirplaneAnim(airplane, time.Second, AirplaneStatusDown, 0)
	}
	c.Add(airplane.CanvasObject)
}

// CreateAirplaneAnim 创建飞机动画
func CreateAirplaneAnim(airplane *Airplane, d time.Duration, airplaneStatus int, repeatCount int) {
	animation := fyne.NewAnimation(d, func(f float32) {
		box := airplane.CanvasObject.(*fyne.Container)
		img := box.Objects[0].(*canvas.Image)
		images := AirplaneImages[airplane.Tp][airplaneStatus]
		l := float32(len(images))
		i := int(f / (1.0 / l))
		var next *canvas.Image
		if i == int(l) {
			next = images[i-1]
		} else {
			next = images[i]
		}
		if img != next {
			box.Remove(img)
			box.Add(next)
		}

		if f == 1 {
			if airplane.Tp == AirplaneHero {
				if airplaneStatus == AirplaneStatusDown {
					// 英雄机死亡，游戏结束
					M.Status = GameStatusEnd
					airplane.Hide()
					M.Pause.Hide()
					M.Menu.Show()
					M.StartButton.Hide()
					M.AgainButton.Show()
					M.EndButton.Show()
					M.Bomb.Hide()
				}
			} else {
				if airplaneStatus == AirplaneStatusDown {
					airplane.ToDefault()
					score, err := M.ScoreData.Get()
					if err != nil {
						fyne.LogError("ScoreData.Get() error", err)
						return
					}
					if err = M.ScoreData.Set(score + airplane.HP); err != nil {
						fyne.LogError("ScoreData.Set() error", err)
					}
				} else if airplaneStatus == AirplaneStatusHit {
					airplane.ToAnim(AirplaneStatusFlying)
				}
			}
		}
	})
	animation.RepeatCount = repeatCount
	animation.Curve = fyne.AnimationLinear
	airplane.anim[airplaneStatus] = animation
}
