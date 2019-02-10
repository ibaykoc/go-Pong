package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ball struct {
	position     pixel.Vec
	velocity     pixel.Vec
	initialSpeed float64
	speed        float64
	radius       float64
	sprite       *pixel.Sprite
}

func newBall(posX, posY, rad, initialSpeed float64) ball {

	pix := make([]color.RGBA, int((rad*2)*(rad*2)))
	c := rad - 0.5
	for y := 0; y < int(rad*2); y++ {
		for x := 0; x < int(rad*2); x++ {
			var index int
			index = y*int(rad*2) + x
			xD := c - float64(x)
			yD := c - float64(y)
			if (xD*xD)+(yD*yD) <= rad*rad {
				pix[index].R = 255
				pix[index].G = 255
				pix[index].B = 255
				pix[index].A = 255
			}
		}
	}
	p := pixel.PictureData{
		Pix:    pix,
		Stride: int(rad) * 2,
		Rect: pixel.Rect{
			Min: pixel.Vec{X: 0, Y: 0},
			Max: pixel.Vec{X: rad * 2, Y: rad * 2},
		},
	}
	s := pixel.NewSprite(&p, p.Bounds())
	return ball{
		position:     pixel.Vec{X: posX, Y: posY},
		velocity:     pixel.Vec{X: initialSpeed, Y: initialSpeed},
		initialSpeed: initialSpeed,
		speed:        initialSpeed,
		radius:       rad,
		sprite:       s,
	}
}

func (b *ball) update(lP, rP *paddle, dt float64) {
	speedMultiplier := 1.1
	b.position.X += b.velocity.X * dt
	b.position.Y += b.velocity.Y * dt

	bMinX := b.position.X - b.radius
	bMaxX := b.position.X + b.radius
	bMinY := b.position.Y - b.radius
	bMaxY := b.position.Y + b.radius

	lPMinX := lP.position.X - lP.size.X/2
	lPMaxX := lP.position.X + lP.size.X/2
	lPMinY := lP.position.Y - lP.size.Y/2
	lPMaxY := lP.position.Y + lP.size.Y/2

	rPMinX := rP.position.X - rP.size.X/2
	rPMaxX := rP.position.X + rP.size.X/2
	rPMinY := rP.position.Y - rP.size.Y/2
	rPMaxY := rP.position.Y + rP.size.Y/2

	// Collide with left paddle
	if bMinX < lPMaxX && bMinX > lPMinX &&
		bMinY <= lPMaxY && bMaxY >= lPMinY {
		b.speed *= speedMultiplier
		dist := (b.position.Y - lP.position.Y) / (lP.size.Y / 2)

		b.position.X = lPMaxX + b.radius
		b.velocity.X = b.speed
		b.velocity.Y = b.speed * dist
	}

	// Collide with right paddle
	if bMaxX > rPMinX && bMaxX < rPMaxX &&
		bMaxY <= rPMaxY && bMinY >= rPMinY {
		b.speed *= speedMultiplier
		dist := (b.position.Y - rP.position.Y) / (rP.size.Y / 2)

		b.position.X = rPMinX - b.radius
		b.velocity.X = -b.speed
		b.velocity.Y = b.speed * dist
	}

	// Collide top
	if bMaxY > wHeight {
		b.position.Y = wHeight - b.radius
		b.velocity.Y = -b.velocity.Y
	}

	// Collide bottom
	if bMinY < 0 {
		b.position.Y = b.radius
		b.velocity.Y = -b.velocity.Y
	}

	// Collide left
	if bMinX < 0 {
		gameStart = false
		rP.score++
		if rP.score >= 3 {
			rP.score = 0
			lP.score = 0
			messageText.Clear()
			changeMessageText([]string{
				rpWinMessage,
				pressToStartMesage},
			)
		}
		b.position.X = wWidth / 2
		b.position.Y = wHeight / 2
		b.speed = b.initialSpeed
		b.velocity = b.velocity.Unit().Scaled(b.speed)

	}

	// Collide right
	if bMaxX > wWidth {
		gameStart = false
		lP.score++
		if lP.score >= 3 {
			rP.score = 0
			lP.score = 0
			messageText.Clear()
			changeMessageText([]string{
				lpWinMessage,
				pressToStartMesage},
			)
		}
		b.position.X = wWidth / 2
		b.position.Y = wHeight / 2
		b.speed = b.initialSpeed
		b.velocity = b.velocity.Unit().Scaled(b.speed)
	}
}

func (b *ball) draw(win *pixelgl.Window) {
	b.sprite.Draw(win, pixel.IM.Moved(b.position))
}
