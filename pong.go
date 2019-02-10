package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	wWidth  = 1024
	wHeight = 768
)

func main() {
	pixelgl.Run(run)
}

var win *pixelgl.Window
var fontAtlas *text.Atlas
var gameStart = false
var messageText *text.Text
var textMat pixel.Matrix

const (
	pressToStartMesage = "Press space to start"
	lpWinMessage       = "You Win!"
	rpWinMessage       = "You Lose!"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "PONG",
		Bounds:    pixel.R(0, 0, wWidth, wHeight),
		VSync:     true,
		Resizable: false,
	}
	var err error
	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	fontAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
	messageText = text.New(pixel.V(0, 0), fontAtlas)
	textMat = pixel.IM
	textMat = textMat.Scaled(messageText.Bounds().Center(), 2)
	textMat = textMat.Moved(pixel.Vec{X: wWidth / 2, Y: wHeight * 0.75})
	changeMessageText([]string{pressToStartMesage})
	b := newBall(win.Bounds().Center().X, win.Bounds().Center().Y, 20, 10)
	lP := newPaddle(20, win.Bounds().Center().Y, 20, 150, 20)
	rP := newPaddle(win.Bounds().W()-20, win.Bounds().Center().Y, 20, 150, 20)
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds() * 60
		last = time.Now()
		win.Clear(colornames.Black)

		lP.update(dt)
		rP.updateAi(&lP, &b, dt)
		if gameStart {
			b.update(&lP, &rP, dt)
		} else {
			messageText.Draw(win, textMat)
			if win.JustReleased(pixelgl.KeySpace) {
				changeMessageText([]string{pressToStartMesage})
				gameStart = true
			}
		}

		lP.draw(win)
		rP.draw(win)
		b.draw(win)

		win.Update()

		if win.JustReleased(pixelgl.KeyEscape) {
			break
		}
	}
}

func changeMessageText(messages []string) {
	messageText.Clear()
	for _, msg := range messages {
		messageText.Dot.X -= messageText.BoundsOf(msg).W() / 2
		fmt.Fprintln(messageText, msg)
	}
}
