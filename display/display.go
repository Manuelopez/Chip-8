package display

import (
	"github.com/nsf/termbox-go"
)

type Display struct {
	Screen       [64][32]bool
	screenWidth  int
	screenHeight int
}

func New() *Display {
	return &Display{screenWidth: 64, screenHeight: 32}
}

func (d *Display) Start() {

	err := termbox.Init()
    
	if err != nil {
		panic(err)
	}


	// Draw the pixels to the screen
	for x := 0; x < d.screenWidth; x++ {
		for y := 0; y < d.screenHeight; y++ {
			termbox.SetCell(x, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
		}
	}

	// Update the screen
    termbox.Flush()

	// Wait for a key press

}

func (d *Display) Update(){
	for x := 0; x < d.screenWidth; x++ {
		for y := 0; y < d.screenHeight; y++ {

            if(d.Screen[x][y]){
			    termbox.SetCell(x, y, ' ', termbox.ColorWhite, termbox.ColorWhite)
            }else{

			    termbox.SetCell(x, y, ' ', termbox.ColorBlack, termbox.ColorBlack)
            }
		}
	}
    termbox.Flush()

}
