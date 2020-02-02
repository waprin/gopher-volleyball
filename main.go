package main

import (
	"fmt"
	"os"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

)

const (
	Width=750
	Height=400
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Printf("could not initialize SDL: %v", err)
		os.Exit(1)
	}
	defer sdl.Quit()


	if err := ttf.Init(); err != nil {
		fmt.Printf("Could not initialize ttf: %v", err)
	}
	window, r, err := sdl.CreateWindowAndRenderer(Width, Height, sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	defer r.Destroy()

	running := true
	g, err := newGame(r, Width, Height)
	if err != nil {
		fmt.Printf("Could not initialize game: %v", err)
	}
	g.start(r)

	for running == true {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch  e := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				switch e.Keysym.Sym {
				case sdl.K_a:
					g.handlePlayer1LeftTouch(e.State)
				case sdl.K_d:
					g.handlePlayer1RightTouch(e.State)
				case sdl.K_w:
					g.handlePlayer1UpTouch(e.State)
				case sdl.K_LEFT:
					g.handlePlayer2LeftTouch(e.State)
				case sdl.K_RIGHT:
					g.handlePlayer2RightTouch(e.State)
				case sdl.K_UP:
					g.handlePlayer2UpTouch(e.State)
				case sdl.K_SPACE:
					g.handleSpaceBar(e.State)
				}
			case *sdl.MouseButtonEvent:
				mouseE := event.(*sdl.MouseButtonEvent)
				g.handleMouseUp(mouseE.X, mouseE.Y)
			}
		}
	}
	return
}

