package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

var BLUE = sdl.Color{0, 0, 255, 255}
var BLACK = sdl.Color{0, 0, 0, 255}
var BACKGROUND_COLOR = BLACK

func main() {
	runtime.LockOSThread()
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		panic(err)
	}
	defer ttf.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	var bgRect sdl.Rect
	surface.GetClipRect(&bgRect)

	font, err := ttf.OpenFont("/Library/Fonts/Arial.ttf", 128)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	fpsTimer := time.NewTicker(time.Second / 60)

MainLoop:
	for {
		select {
		case <-fpsTimer.C:
			text, err := font.RenderUTF8_Solid(fmt.Sprintf("%d", rand.Int31n(10)), BLUE)
			if err != nil {
				panic(err)
			}

			surface.FillRect(&bgRect, BACKGROUND_COLOR.Uint32())

			// rect := sdl.Rect{0, 0, 200, 200}
			// surface.FillRect(&rect, 0xffff0000)
			var rect sdl.Rect
			text.GetClipRect(&rect)
			rect.Y += 300 - (rect.H / 2)
			text.Blit(nil, surface, &rect)

			window.UpdateSurface()
		default:
			for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
				switch t := e.(type) {
				case *sdl.QuitEvent:
					break MainLoop
				case *sdl.KeyDownEvent:
					if t.Keysym.Sym == sdl.K_ESCAPE {
						break MainLoop
					}
				}
			}
		}

	}

	fmt.Println("main thread exiting")
}
