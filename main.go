package main

import (
	"math"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const (
	GameStatusDefault = iota // 初始状态
	GameStatusGaming         // 游戏中
	GameStatusPause          // 暂停中
	GameStatusEnd            // 游戏结束
)

func GetBackground() fyne.CanvasObject {
	image := canvas.NewImageFromResource(resourceBackgroundPng)
	image.Resize(fyne.NewSize(MasterWidth, MasterHeight))

	image1 := canvas.NewImageFromResource(resourceBackgroundPng)
	image1.Resize(fyne.NewSize(MasterWidth, MasterHeight))
	image1.Move(fyne.NewPos(0, MasterHeight))

	bg := container.NewWithoutLayout(image, image1)
	bg.Resize(fyne.NewSize(MasterWidth, MasterHeight*2))

	bgAni := canvas.NewPositionAnimation(
		fyne.NewPos(0, -MasterHeight),
		fyne.NewPos(0, 0),
		time.Second*20,
		func(position fyne.Position) {
			bg.Move(position)
			bg.Refresh()
		})
	bgAni.RepeatCount = fyne.AnimationRepeatForever
	bgAni.Curve = fyne.AnimationLinear
	bgAni.Start()
	return bg
}

type Manager struct {
	Pause          fyne.CanvasObject
	Score          fyne.CanvasObject
	Bomb           fyne.CanvasObject
	Menu           fyne.CanvasObject
	StartButton    fyne.CanvasObject
	AgainButton    fyne.CanvasObject
	EndButton      fyne.CanvasObject
	Hero           fyne.CanvasObject
	ScoreData      binding.Int              // 分数
	Bullets        [BulletMax][]*Bullet     // 所有子弹
	Airplanes      [AirplaneMax][]*Airplane // 所有飞机
	Props          [PropMax][]*Prop         // 所有道具
	Status         int                      // 游戏状态
	LastEnemyTime  time.Time                // 最后产生敌机的时间
	LastPropTime   time.Time                // 最后产生道具的时间
	LastBulletTime time.Time                // 最后发射子弹的时间
	HeroClick      bool                     // 英雄机是否被左键按下
	SubPos         fyne.Position            // 记录按下左键的位置和英雄机位置的差值
	DoubleNum      int                      // 双倍子弹剩余数量
	BombNum        int                      // 剩余炸弹数量
}

var M = &Manager{
	Bullets:   [BulletMax][]*Bullet{},
	Airplanes: [AirplaneMax][]*Airplane{},
	Props:     [PropMax][]*Prop{},
}

var MyApp = app.New()

var Ch chan func()

func Done(f func()) {
	select {
	case Ch <- f:
	default:
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
	Ch = make(chan func(), 200)
	// 子弹
	for i := 0; i < 100; i++ {
		M.Bullets[BulletHero] = append(M.Bullets[BulletHero], NewBullet(BulletCreateHero, BulletStatusUnused))
		//M.Bullets[BulletEnemy] = append(M.Bullets[BulletEnemy], NewBullet(BulletCreateEnemy, BulletStatusUnused))
	}

	// 飞机
	for i := 0; i < 50; i++ {
		M.Airplanes[AirplaneEnemy1] = append(M.Airplanes[AirplaneEnemy1], NewAirplane(AirplaneEnemy1, AirplaneStatusDefault))
	}
	for i := 0; i < 5; i++ {
		M.Airplanes[AirplaneEnemy2] = append(M.Airplanes[AirplaneEnemy2], NewAirplane(AirplaneEnemy2, AirplaneStatusDefault))
	}
	for i := 0; i < 1; i++ {
		M.Airplanes[AirplaneEnemy3] = append(M.Airplanes[AirplaneEnemy3], NewAirplane(AirplaneEnemy3, AirplaneStatusDefault))
	}
	M.Airplanes[AirplaneHero] = []*Airplane{NewAirplane(AirplaneHero, AirplaneStatusDefault)}

	// 道具
	for i := 0; i < 2; i++ {
		M.Props[PropBomb] = append(M.Props[PropBomb], NewProp(PropBomb, PropStatusDefault))
	}
	for i := 0; i < 4; i++ {
		M.Props[PropBullet] = append(M.Props[PropBullet], NewProp(PropBullet, PropStatusDefault))
	}
	// 分数
	M.ScoreData = binding.NewInt()
}

func main() {
	w := MyApp.NewWindow("飞机大战")
	w.Resize(fyne.NewSize(MasterWidth, MasterHeight))
	w.SetPadded(false)
	w.SetFixedSize(true)
	w.SetMaster()
	w.SetIcon(resourceLifePng)

	// 绘制所有游戏对象
	Draw(w)

	// 状态初始化
	GameStatusDefaultDone()

	// 接收状态修改或处理游戏逻辑
	go Run()

	w.ShowAndRun()
}

func Draw(w fyne.Window) {
	c := container.NewWithoutLayout()
	w.SetContent(c)

	// 监听鼠标位置
	c.Add(NewBoxMouse(fyne.NewSize(MasterWidth, MasterHeight)))

	// 背景图
	c.Add(GetBackground())

	// 子弹
	for _, v := range M.Bullets {
		for _, vv := range v {
			AddBullet(c, vv)
		}
	}

	// 飞机
	for _, v := range M.Airplanes {
		for _, vv := range v {
			AddAirplane(c, vv)
		}
	}
	M.Hero = M.Airplanes[AirplaneHero][0]

	// 道具
	for _, v := range M.Props {
		for _, vv := range v {
			AddProp(c, vv)
		}
	}

	// 暂停按钮
	pause := NewPauseButton(func() {}, func() {
		Done(func() {
			switch M.Status {
			case GameStatusGaming:
				M.Status = GameStatusPause
				M.Menu.Show()
				M.StartButton.Hide()
				M.AgainButton.Show()
				M.EndButton.Show()
			case GameStatusPause:
				M.Status = GameStatusGaming
				M.Menu.Hide()
			}
		})
	})
	M.Pause = pause
	c.Add(pause)

	// 分数
	score := widget.NewLabelWithData(binding.IntToStringWithFormat(M.ScoreData, "%d"))
	score.TextStyle.Bold = true
	score.Resize(fyne.NewSize(MasterWidth, 45))
	score.Move(fyne.NewPos(64, 0))
	M.Score = score
	c.Add(score)

	// 可用道具
	var bomb *BombButton
	bomb = NewBombButton(func() {
		Done(func() {
			if M.Status != GameStatusGaming {
				return
			}

			// 所有敌机死亡
			for _, v := range M.Airplanes[:AirplaneHero] {
				for _, vv := range v {
					if vv.Status != AirplaneStatusFlying {
						continue
					}
					vv.Status = AirplaneStatusDown
					vv.ToAnim(AirplaneStatusDown)
				}
			}

			M.BombNum--
			if M.BombNum <= 0 {
				M.Bomb.Hide()
			}
		})
	})
	M.Bomb = bomb
	c.Add(bomb)

	startFunc := func() {
		M.StartButton.Hide()
		M.Pause.Show()
		M.Score.Show()
		M.Status = GameStatusGaming
		M.Airplanes[AirplaneHero][0].Status = AirplaneStatusFlying
		M.LastPropTime = time.Now().Add(time.Second * 5)
	}

	// 菜单
	M.StartButton = NewItemButton(resourceGamestartPng, func() {
		// 开始游戏
		Done(startFunc)
	})
	M.AgainButton = NewItemButton(resourceAgainPng, func() {
		// 重新开始
		Done(func() {
			GameStatusDefaultDone()
			startFunc()
		})
	})
	M.EndButton = NewItemButton(resourceGameoverPng, func() {
		// 结束游戏
		Done(func() {
			GameStatusDefaultDone()
			M.Status = GameStatusDefault
		})
	})
	menu := container.NewCenter(container.NewVBox(M.StartButton, M.AgainButton, M.EndButton))
	menu.Resize(fyne.NewSize(MasterWidth, MasterHeight))
	M.Menu = menu
	c.Add(menu)
}

func Run() {
	for range time.Tick(time.Second / time.Duration(GameSpeed)) {
		select {
		case f := <-Ch:
			f()
		default:
		}

		// 游戏状态处理
		switch M.Status {
		case GameStatusDefault:

		case GameStatusGaming:
			GameStatusGamingDone()

		case GameStatusPause:

		case GameStatusEnd:

		}
	}
}

func GameStatusDefaultDone() {
	M.LastEnemyTime = time.Time{}
	M.LastPropTime = time.Time{}
	M.LastBulletTime = time.Time{}
	M.Pause.Hide()
	M.Pause.(*PauseButton).ToDefault()

	M.Score.Hide()
	if err := M.ScoreData.Set(0); err != nil {
		fyne.LogError("ScoreData.Set(0) error", err)
	}

	M.Bomb.Hide()
	M.BombNum = 0

	M.Menu.Show()
	M.StartButton.Show()
	M.AgainButton.Hide()
	M.EndButton.Hide()
	M.Hero.Show()

	for _, v := range M.Airplanes {
		for _, vv := range v {
			vv.ToDefault()
		}
	}

	for _, v := range M.Bullets {
		for _, vv := range v {
			vv.ToDefault()
		}
	}

	for _, v := range M.Props {
		for _, vv := range v {
			vv.ToDefault()
		}
	}

	M.HeroClick = false
}

func GameStatusGamingDone() {
	curTime := time.Now()
	// 新增敌机
	if curTime.Sub(M.LastEnemyTime) > time.Second/AirplaneHz {
		var has bool
		for i := 0; i < 10; i++ {
			for _, v := range M.Airplanes[rand.Intn(AirplaneHero)] {
				if v.Status != AirplaneStatusDefault {
					continue
				}
				has = true

				// 敌机入场状态
				v.Status = AirplaneStatusFlying
				switch v.Tp {
				case AirplaneEnemy1:
					v.speed = float32(Enemy2FlyingSpeedMax+rand.Intn(Enemy1FlyingSpeedMax-Enemy2FlyingSpeedMax)) /
						float32(GameSpeed)
				case AirplaneEnemy2:
					v.speed = float32(Enemy3FlyingSpeedMax+rand.Intn(Enemy2FlyingSpeedMax-Enemy3FlyingSpeedMax)) /
						float32(GameSpeed)
				case AirplaneEnemy3:
					v.speed = float32(EnemyFlyingSpeedMin+rand.Intn(Enemy3FlyingSpeedMax-EnemyFlyingSpeedMin)) /
						float32(GameSpeed)
				}
				v.ToAnim(v.Status)
				v.Move(fyne.NewPos(float32(rand.Intn(int(MasterWidth-v.Size().Width))), v.Position().Y))
				v.lastBulletTime = time.Time{}
				//log.Println("add airplane", v.Position())

				break
			}
			if !has {
				continue
			}
			M.LastEnemyTime = curTime
			break
		}
		if !has {
			M.LastEnemyTime = curTime.Add(time.Second)
			//log.Println("可能没有可用的敌机了")
		}
	}

	// 新增道具
	if curTime.Sub(M.LastPropTime) > PropHz*time.Second {
		var has bool
		for i := 0; i < 6; i++ {
			for _, v := range M.Props[rand.Intn(PropMax)] {
				if v.Status != PropStatusDefault {
					continue
				}
				has = true

				// 道具入场状态
				v.Status = PropStatusFlying
				v.speed = float32(PropSpeed) / float32(GameSpeed)
				v.Move(fyne.NewPos(float32(rand.Intn(int(MasterWidth-v.Size().Width))), v.Position().Y))

				break
			}
			if !has {
				continue
			}
			M.LastPropTime = curTime
			break
		}
		if !has {
			M.LastPropTime = curTime.Add(time.Second)
			//log.Println("可能没有可用的道具了")
		}
	}

	// 新增子弹
	if curTime.Sub(M.LastBulletTime) > time.Second/BulletHz {
		var bullets []*Bullet
		for _, v := range M.Bullets[BulletHero] {
			if v.Status == BulletStatusUnused {
				bullets = append(bullets, v)
			}
		}
		for {
			if len(bullets) == 0 {
				fyne.LogError("bullet not enough", nil)
				break
			}
			if M.DoubleNum > 0 && len(bullets) < 2 {
				fyne.LogError("bullet not enough", nil)
				break
			}

			f := func(b *Bullet, pos fyne.Position, res fyne.Resource) {
				b.Status = BulletStatusFlying
				if res != nil {
					b.CanvasObject.(*canvas.Image).Resource = res
					b.Refresh()
				}
				b.Move(pos)
			}

			pos := M.Hero.Position().Add(fyne.NewPos(M.Hero.Size().Width/2, M.Hero.Size().Height/4))
			if M.DoubleNum > 0 {
				M.DoubleNum -= 2
				f(bullets[0], pos.Subtract(fyne.NewPos(M.Hero.Size().Width/4+bullets[0].Size().Width, 0)), resourceBullet1Png)
				f(bullets[1], pos.Add(fyne.NewPos(M.Hero.Size().Width/4, 0)), resourceBullet1Png)
			} else {
				f(bullets[0], pos.Subtract(fyne.NewPos(1, 0)), resourceBullet2Png)
			}
			break
		}
		M.LastBulletTime = curTime
	}

	// 子弹移动
	for _, v := range M.Bullets {
		for _, vv := range v {
			if vv.Status == BulletStatusUnused {
				continue
			}
			if vv.Tp == BulletCreateEnemy {
				vv.Move(vv.Position().Add(fyne.NewPos(0, float32(BulletSpeed)/GameSpeed)))
				if vv.Position().Y > MasterHeight {
					vv.ToDefault()
				}
			} else {
				vv.Move(vv.Position().Add(fyne.NewPos(0, -float32(BulletSpeed)/GameSpeed)))
				if vv.Position().Y < -vv.Size().Height {
					vv.ToDefault()
				}
			}
		}
	}

	// 敌机移动
	for _, v := range M.Airplanes[:AirplaneHero] {
		for _, vv := range v {
			if vv.Status == AirplaneStatusDefault || vv.Status == AirplaneStatusDown {
				continue
			}
			vv.Move(vv.Position().Add(fyne.NewPos(0, vv.speed)))
			if vv.Position().Y > MasterHeight {
				vv.ToDefault()
			}
		}
	}

	// 道具移动
	for _, v := range M.Props {
		for _, vv := range v {
			if vv.Status == PropStatusDefault {
				continue
			}
			vv.Move(vv.Position().Add(fyne.NewPos(0, vv.speed)))
			if vv.Position().Y > MasterHeight {
				vv.ToDefault()
			}
		}
	}

	// 碰撞检测
	// 子弹和敌机
	for _, v := range M.Bullets[BulletHero] {
		for _, vv := range M.Airplanes[:AirplaneHero] {
			for _, vvv := range vv {
				if CrashBullet(v, vvv) {
					// 子弹
					v.ToDefault()
					// 敌机,击中动画，扣血量
					vvv.ToHit()
				}
			}
		}
	}
	// 英雄和道具
	for _, v := range M.Props {
		for _, vv := range v {
			if CrashProp(M.Airplanes[AirplaneHero][0], vv) {
				vv.ToDefault()
				switch vv.Tp {
				case PropBullet:
					// 50发双倍子弹
					M.DoubleNum += 50
				case PropBomb:
					if M.BombNum <= 0 {
						M.Bomb.Show()
					}
					M.BombNum++
				}
			}
		}
	}

	// 英雄和敌机
	for _, v := range M.Airplanes[:AirplaneHero] {
		for _, vv := range v {
			if CrashEnemy(M.Airplanes[AirplaneHero][0], vv) {
				vv.ToHit()
				M.Airplanes[AirplaneHero][0].ToHit()
			}
		}
	}
}

func distance(a, b fyne.CanvasObject) bool {
	pos := a.Position().Add(fyne.NewPos(a.Size().Width/2, a.Size().Height/2))
	pos2 := b.Position().Add(fyne.NewPos(b.Size().Width/2, b.Size().Height/2))

	a1 := math.Abs(float64(pos.X - pos2.X))
	b1 := math.Abs(float64(pos.Y - pos2.Y))

	return math.Sqrt(a1*a1+b1*b1) < float64(a.Size().Width/3+b.Size().Width/2)
}

func CrashBullet(bullet *Bullet, enemy *Airplane) bool {
	if bullet.Status != BulletStatusFlying {
		return false
	}
	if enemy.Status != AirplaneStatusFlying && enemy.Status != AirplaneStatusHit {
		return false
	}
	return distance(bullet, enemy)
}

func CrashProp(hero *Airplane, prop *Prop) bool {
	if hero.Status != AirplaneStatusFlying {
		return false
	}
	if prop.Status != PropStatusFlying {
		return false
	}
	return distance(hero, prop)
}

func CrashEnemy(hero, enemy *Airplane) bool {
	if hero.Status != AirplaneStatusFlying {
		return false
	}
	if enemy.Status != AirplaneStatusFlying {
		return false
	}
	return distance(hero, enemy)
}
