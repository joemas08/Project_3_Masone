package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/png"
	"log"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

const (
	GameWidth   = 1400
	GameHeight  = 700
	PlayerSpeed = 10
)

type Sprite struct {
	pict *ebiten.Image
	xloc int
	yloc int
	dX   int
	dY   int
}

type enemySprite struct {
	pict *ebiten.Image
	xloc int
	yloc int
}

type Game struct {
	player   Sprite
	enemy    enemySprite
	score    int
	drawOps  ebiten.DrawImageOptions
	gameName string
}

func (g *Game) Update() error {
	processPlayerInput(g)
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	g.drawOps.GeoM.Reset()
	g.drawOps.GeoM.Translate(float64(g.enemy.xloc), float64(g.enemy.yloc))
	screen.DrawImage(g.enemy.pict, &g.drawOps)
	g.drawOps.GeoM.Reset()
	g.drawOps.GeoM.Translate(float64(g.player.xloc), float64(g.player.yloc))
	screen.DrawImage(g.player.pict, &g.drawOps)

}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("assets")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("assets/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func processPlayerInput(theGame *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		theGame.player.dY = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		theGame.player.dX = PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		theGame.player.dX = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		theGame.player.dY = PlayerSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		theGame.player.dX = 0
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		theGame.player.dY = 0
	}
	theGame.player.yloc += theGame.player.dY
	theGame.player.xloc += theGame.player.dX
	if theGame.player.yloc <= 0 {
		theGame.player.dY = 0
		theGame.player.yloc = 0
	} else if theGame.player.yloc+theGame.player.pict.Bounds().Size().Y > GameHeight {
		theGame.player.dY = 0
		theGame.player.yloc = GameHeight - theGame.player.pict.Bounds().Size().Y
	}
	if theGame.player.xloc <= 0 {
		theGame.player.dX = 0
		theGame.player.xloc = 0
	} else if theGame.player.xloc+theGame.player.pict.Bounds().Size().X > GameWidth {
		theGame.player.dX = 0
		theGame.player.xloc = GameWidth - theGame.player.pict.Bounds().Size().X
	}
}

func main() {
	carGame := Game{score: 0, gameName: "Car Game"}
	ebiten.SetWindowSize(GameWidth, GameHeight)
	ebiten.SetWindowTitle(carGame.gameName)
	carGame.player = Sprite{
		pict: loadPNGImageFromEmbedded("car.png"),
		xloc: 200,
		yloc: 300,
		dX:   0,
		dY:   0,
	}

	carGame.enemy = enemySprite{
		pict: loadPNGImageFromEmbedded("jerry_can.png"),
		xloc: 300,
		yloc: 300,
	}
	if err := ebiten.RunGame(&carGame); err != nil {
		log.Fatal("Oh no! something terrible happened and the game crashed", err)
	}
}
