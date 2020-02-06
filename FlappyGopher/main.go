package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"time"
)

func main() {
	fmt.Println("Hey")

	test()
	//if err := run(); err != nil {
	//	fmt.Fprintf(os.Stderr, "%v", err)
	//	os.Exit(2)
	//}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize %v", err)
	}
	defer sdl.Quit()

	ttf.Init()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()
	w.Maximize()
	_ = r

	x := sdl.PollEvent()
	quit := false
	for !quit {
		if x.GetType() == sdl.QUIT {
			quit = true
		}
		fmt.Println("Waiting....")
	}

	time.Sleep(5 * time.Second)

	return nil
}

func test() {
	// try to initialize everything
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize sdl: %s\n", err)
		os.Exit(1)
	}

	// try to create a window
	window, err := sdl.CreateWindow("Go + SDL2 Lesson 1", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		640, 480, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(2)
	}

	// window has been created, now need to get the window surface to draw on window
	screenSurface, err := window.GetSurface()
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create surface: %s\n", err)
		os.Exit(2)
	}

	// create the first rectangle (x position, y position, width, height)
	rect := sdl.Rect{300, 300, 100, 80}
	// draw the rect on the window surface and choose color based on r,g,b - in this case the color blue
	screenSurface.FillRect(&rect, sdl.MapRGB(screenSurface.Format, 0, 0, 255))

	// draw second rectangle (this one green) - demonstrates the fact your can create rects as needed inside function calls
	//
	// the first argument in FillRect wants a pointer and neither rect we have used was a pointer so the & was used in both cases
	// to provide the needed pointer
	screenSurface.FillRect(&sdl.Rect{0, 0, 200, 200}, sdl.MapRGB(screenSurface.Format, 0, 255, 0))

	// if nil is used as the first argument instead of a rect that tells sdl to draw the rect on the entire window surface area
	// uncomment the next line to see the entire window in red
	//screenSurface.FillRect(nil, sdl.MapRGB(screenSurface.Format, 255, 0, 0))

	// it is not enough to draw on the window surface, you must tell sdl to show what you've done
	window.UpdateSurface()

	// used to keep window open for five seconds
	time.Sleep(time.Second * 5)

	// program is over, time to start shutting down. Keep in mind that sdl is written in C and does not have convenient
	// garbage collection like Go does
	window.Destroy()

	sdl.Quit()
}
