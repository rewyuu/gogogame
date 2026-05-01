package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GopherAnim struct {
	frames  []*ebiten.Image
	current int
	timer   float64
	fps     float64
	x       float64
	y       float64
	scale   float64
}

func NewGopherAnim(folder string, frameCount int, x, y float64, scale float64, startFrame int) *GopherAnim {
	frames := make([]*ebiten.Image, frameCount)
	for i := range frames {
		path := fmt.Sprintf("assets/spritesheet/%s/frame%04d.png", folder, i+1)
		frames[i] = loadImage(path)
	}
	return &GopherAnim{
		frames:  frames,
		fps:     12,
		x:       x,
		y:       y,
		scale:   scale,
		current: startFrame,
	}
}

func (a *GopherAnim) Update() {
	a.timer += 1.0 / 60.0
	if a.timer >= 1.0/a.fps {
		a.timer = 0
		a.current++
		if a.current >= len(a.frames) {
			a.current = 0
		}
	}
}

func (a *GopherAnim) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(a.scale, a.scale)
	op.GeoM.Translate(a.x, a.y)
	screen.DrawImage(a.frames[a.current], op)
}

// scene particles — reusing same Particle struct from main.go
type SceneParticle struct {
	x, y, speed, phase, drift, size float64
}

func spawnSceneParticles(count int) []*SceneParticle {
	particles := make([]*SceneParticle, count)
	for i := range particles {
		particles[i] = &SceneParticle{
			x:     float64(rand.Intn(320)),
			y:     float64(rand.Intn(240)),
			speed: 0.2 + rand.Float64()*0.4,
			phase: rand.Float64() * math.Pi * 2,
			drift: 0.3 + rand.Float64()*0.8,
			size:  1 + rand.Float64()*2,
		}
	}
	return particles
}

type GameScene struct {
	gophers   []*GopherAnim
	particles []*SceneParticle
	time      float64
	titleY    float64   // for float animation
}

func NewGameScene() *GameScene {
	return &GameScene{
		gophers: []*GopherAnim{
			//path, frame count, x, y, scale, start frame
			// dancers
			NewGopherAnim("gopher_dance",  24, 10, 190, 0.5, 0),
			NewGopherAnim("gopher_dance",  24, 35, 195, 0.45, 2),
			NewGopherAnim("gopher_dance",  24, 55, 200, 0.55, 4),
			NewGopherAnim("gopher_dance",  24, 80, 215, 0.3, 6),
			NewGopherAnim("gopher_dance",  24, 90, 195, 0.7, 8),
			NewGopherAnim("gopher_dance",  24, 120, 175, 1, 10),
			NewGopherAnim("gopher_dance",  24, 165, 200, 0.5, 12),
			NewGopherAnim("gopher_dance",  24, 190, 225, 0.2, 14),
			NewGopherAnim("gopher_dance",  24, 200, 225, 0.2, 16),
			NewGopherAnim("gopher_dance",  24, 210, 225, 0.2, 18),
			NewGopherAnim("gopher_dance",  24, 225, 210, 0.45, 20),
			NewGopherAnim("gopher_dance",  24, 255, 220, 0.3, 22),
			NewGopherAnim("gopher_dance",  24, 270, 195, 0.7, 23),

			NewGopherAnim("gopher_dance",  24, 150, 100, 0.3, 6),
			NewGopherAnim("gopher_dance",  24, 100, 100, 0.7, 8),

			NewGopherAnim("gopher_dance",  24, 95, 50, 0.3, 4),
			NewGopherAnim("gopher_dance",  24, 115, 50, 0.3, 8),
			NewGopherAnim("gopher_dance",  24, 135, 50, 0.3, 12),
			NewGopherAnim("gopher_dance",  24, 155, 50, 0.3, 16),
			NewGopherAnim("gopher_dance",  24, 175, 50, 0.3, 20),
			NewGopherAnim("gopher_dance",  24, 195, 50, 0.3, 23),

			//beer
			NewGopherAnim("gopher_beer",   10, 35,  110, 0.5, 0),
			NewGopherAnim("gopher_beer",   10, 55,  125, 0.7, 2),
			NewGopherAnim("gopher_beer",   10, 20,  140, 0.3, 4),
			NewGopherAnim("gopher_beer",   10, 15,  85, 0.4, 6),


			//coffee
			NewGopherAnim("gopher_coffee", 10, 225, 175, 0.3, 0),
			NewGopherAnim("gopher_coffee", 10, 260, 160, 0.3, 2),
			NewGopherAnim("gopher_coffee", 10, 75, 75, 0.3, 4),
			NewGopherAnim("gopher_coffee", 10, 50, 50, 0.3, 6),

			NewGopherAnim("gopher_coffee", 10, 250, 50, 0.8, 8),
			NewGopherAnim("gopher_coffee", 10, 210, 100, 0.8, 5),
			NewGopherAnim("gopher_coffee", 10, 170, 125, 0.8, 3),


			//ninja
			NewGopherAnim("gopher_ninja",  8,  275, 10, 0.5, 0),
			NewGopherAnim("gopher_ninja",  8,  0, 10, 0.8, 5),
		},
		particles: spawnSceneParticles(60),
		titleY:    30,
	}
}

func (gs *GameScene) Update(g *Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = StateMenu
		g.fadeAlpha = 0
		g.musicPlayer.Pause()
	}

	gs.time += 0.016

	for _, gopher := range gs.gophers {
		gopher.Update()
	}

	// update particles
	for _, p := range gs.particles {
		p.y -= p.speed
		p.x += math.Sin(gs.time+p.phase) * p.drift
		if p.y < -10 {
			p.y = 250
			p.x = float64(rand.Intn(320))
		}
	}

	return nil
}

func (gs *GameScene) Draw(screen *ebiten.Image, fontFace text.Face) {
	screen.Fill(color.RGBA{10, 10, 30, 255})

	// particles
	for _, p := range gs.particles {
		ebitenutil.DrawRect(screen, p.x, p.y, p.size, p.size, color.RGBA{255, 255, 255, 60})
	}

	// gophers
	for _, gopher := range gs.gophers {
		gopher.Draw(screen)
	}

	// pastel title — floating
	titleFloat := gs.titleY + math.Sin(gs.time*1.5)*3

	t := gs.time * 5
	r  := float32(0.85 + math.Sin(t)*0.15)
	gr := float32(0.7  + math.Sin(t+2.1)*0.15)
	b  := float32(0.9  + math.Sin(t+4.2)*0.1)

	op := &text.DrawOptions{}
	op.GeoM.Translate(100, titleFloat)
	op.ColorScale.Scale(r, gr, b, 1)
	text.Draw(screen, "GOPHER FUN PARTY", fontFace, op)
}