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

var font *ttf.Font

type Num struct {
	n   int32
	s   *sdl.Surface
	r   sdl.Rect
	pos int32
}

func NewNum() *Num {
	num := &Num{n: rand.Int31n(10)}
	s, err := font.RenderUTF8_Solid(fmt.Sprintf("%d", num.n), BLUE)
	if err != nil {
		panic(err)
	}
	num.s = s
	num.s.GetClipRect(&num.r)
	num.r.Y += 300 - (num.r.H / 2)
	go num.Think()
	return num
}

func (num *Num) Think() {
	ticker := time.NewTicker(time.Second / 10)
	for {
		select {
		case <-ticker.C:
			num.pos += 10
		}
	}
}

func (num *Num) Draw(surface *sdl.Surface) {
	rect := num.r
	rect.X += num.pos
	num.s.Blit(nil, surface, &rect)
}

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

	font, err = ttf.OpenFont("/Library/Fonts/Arial.ttf", 128)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	fpsTimer := time.NewTicker(time.Second / 60)

	newNumTicker := time.NewTicker(time.Second)

	nums := make(map[*Num]*Num)

	go func() {
		for {
			<-newNumTicker.C
			num := NewNum()
			nums[num] = num
		}
	}()

MainLoop:
	for {
		select {
		case <-fpsTimer.C:

			surface.FillRect(&bgRect, BACKGROUND_COLOR.Uint32())

			for _, num := range nums {
				num.Draw(surface)
			}

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
