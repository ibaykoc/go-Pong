package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type paddle struct {
	position pixel.Vec
	velocity pixel.Vec
	size     pixel.Vec
	sprite   *pixel.Sprite
	text     *text.Text
	score    int
	speed    float64
}

func newPaddle(posX, posY, sizeX, sizeY, speed float64) paddle {
	pix := make([]color.RGBA, int((sizeX)*(sizeY)))
	for y := 0; y < int(sizeY); y++ {
		for x := 0; x < int(sizeX); x++ {
			index := y*int(sizeX) + x
			pix[index].R = 255
			pix[index].G = 255
			pix[index].B = 255
			pix[index].A = 255
		}
	}
	p := pixel.PictureData{
		Pix:    pix,
		Stride: int(sizeX),
		Rect: pixel.Rect{
			Min: pixel.Vec{X: 0, Y: 0},
			Max: pixel.Vec{X: sizeX, Y: sizeY},
		},
	}
	s := pixel.NewSprite(&p, p.Bounds())
	basicTxt := text.New(pixel.V(0, 0), fontAtlas)

	return paddle{
		position: pixel.Vec{X: posX, Y: posY},
		velocity: pixel.Vec{X: 0, Y: 0},
		size:     pixel.Vec{X: sizeX, Y: sizeY},
		sprite:   s,
		text:     basicTxt,
		score:    0,
		speed:    speed,
	}
}

func (p *paddle) update(dt float64) {
	p.velocity.Y = 0

	pMinY := p.position.Y - p.size.Y/2
	pMaxY := p.position.Y + p.size.Y/2

	if win.Pressed(pixelgl.KeyDown) && pMinY > 0 {
		p.velocity.Y -= p.speed
	}
	if win.Pressed(pixelgl.KeyUp) && pMaxY < wHeight {
		p.velocity.Y += p.speed
	}

	p.position.X += p.velocity.X * dt
	p.position.Y += p.velocity.Y * dt
}

func (p *paddle) updateAi(lP *paddle, b *ball, dt float64) {
	p.velocity.Y = 0
	pTop := p.position.Y + p.size.Y/2
	pBot := p.position.Y - p.size.Y/2

	if lP.position.Y > wWidth/2 {
		// if lP on top
		// bounce the ball with bottom part
		pMidBot := p.position.Y - p.size.Y/4

		if b.position.Y < pMidBot {
			// if ball below
			// Move down
			p.velocity.Y = -p.speed
		} else if b.position.Y > pMidBot {
			// if ball above
			// Move up
			p.velocity.Y = p.speed
		}
	} else {
		// if lP on middle or bottom
		// bounce the ball with top part
		pMidTop := p.position.Y + p.size.Y/4

		if b.position.Y > pMidTop {
			// if ball above
			// Move up
			p.velocity.Y = p.speed
		} else if b.position.Y < pMidTop {
			// if ball below
			// Move down
			p.velocity.Y = -p.speed
		}
	}

	p.position.X += p.velocity.X * dt
	p.position.Y += p.velocity.Y * dt

	topDelta := wHeight - pTop
	if topDelta < 0 {
		p.position.Y += topDelta
	}
	if pBot < 0 {
		p.position.Y -= pBot
	}
}

func (p *paddle) draw(win *pixelgl.Window) {
	p.text.Clear()
	fmt.Fprintf(p.text, "%d", p.score)
	p.sprite.Draw(win, pixel.IM.Moved(p.position))
	tMat := pixel.IM
	tMat = tMat.Scaled(p.text.Bounds().Center(), 4)
	tMat = tMat.Moved(pixel.Vec{X: lerp(0.5, p.position.X, wWidth/2), Y: wHeight - 30})
	p.text.Draw(win, tMat)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}
