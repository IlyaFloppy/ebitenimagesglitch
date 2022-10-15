package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	game := Game{
		renderedFirstFrame: make(chan struct{}),
	}

	go func() {
		<-game.renderedFirstFrame

		for key := range dataByKey {
			getImage(key)
		}

		imagesCacheMx.Lock()
		defer imagesCacheMx.Unlock()
		for key, img := range imagesCache {
			f, err := os.Create(fmt.Sprintf("./out/%s.png", key))
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()

			err = png.Encode(f, img)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("saved %s\t/\t%d", key, len(imagesCache))
		}

		os.Exit(0)
	}()

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

var imagesCache = map[string]*ebiten.Image{}
var imagesCacheMx sync.Mutex

// readonly
var dataByKey = map[string][]byte{
	"image1": image1,
	"image2": image2,
	"image3": image3,
	"image4": image4,
	"image5": image5,
	"image6": image6,
	"image7": image7,
	"image8": image8,
	"image9": image9,
}

func getImage(key string) *ebiten.Image {
	imagesCacheMx.Lock()
	defer imagesCacheMx.Unlock()

	if img, ok := imagesCache[key]; ok {
		return img
	}

	source, _, err := image.Decode(bytes.NewBuffer(dataByKey[key]))
	if err != nil {
		log.Fatalln(err)
	}

	img := ebiten.NewImageFromImageWithOptions(source, &ebiten.NewImageFromImageOptions{
		Unmanaged: true,
	})

	imagesCache[key] = img

	return img
}

type Game struct {
	renderedFirstFrame chan struct{}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{
		R: 255,
		G: 50,
		B: 50,
		A: 255,
	})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("fps=%0.2f\ntps=%0.2f", ebiten.CurrentFPS(), ebiten.CurrentTPS()))

	select {
	case <-g.renderedFirstFrame:
	default:
		close(g.renderedFirstFrame)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

//go:embed images/1.png
var image1 []byte

//go:embed images/2.png
var image2 []byte

//go:embed images/3.png
var image3 []byte

//go:embed images/4.png
var image4 []byte

//go:embed images/5.png
var image5 []byte

//go:embed images/6.png
var image6 []byte

//go:embed images/7.png
var image7 []byte

//go:embed images/8.png
var image8 []byte

//go:embed images/9.png
var image9 []byte
