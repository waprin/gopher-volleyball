package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Fprint(os.Stderr, "could not initialize SDL: %v", err)
		os.Exit(1)
	}
	defer sdl.Quit()


	if err := ttf.Init(); err != nil {
		fmt.Printf("Could not initialize ttf: %v", err)
	}
	window, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	defer r.Destroy()

	running := true
	g := newGame()
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
					fmt.Printf("the A button was pressed \n")
					g.handleLeftTouch()
				case sdl.K_d:
					fmt.Printf("The d button was pressed\n")
				}
			}
		}
	}



	fmt.Printf("done main thread\n")
	return
}

