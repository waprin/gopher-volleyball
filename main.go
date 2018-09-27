package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
	"os"
	"github.com/veandco/go-sdl2/ttf"
	"time"
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

	g := &game{leftX:50}


	running := true

	lastUpdate = time.Now()
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
			default:
				diff := time.Since(lastUpdate)
				if diff > 30 * time.Millisecond {
					fmt.Printf("Diff is %v\n", diff)
					lastUpdate = time.Now()
					g.tick()
					g.render(r)
				}
			}
		}
	}
	fmt.Printf("done main thread\n")
	return
}

