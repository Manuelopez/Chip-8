package main

import (
	"chip-8/display"
	"time"
)

const (
    screenWidth  = 64
    screenHeight = 32
)

func main() {
    d := display.New()
    d.Start()
    time.Sleep(time.Second*5)
    d.Update()
    time.Sleep(time.Second*5)
}

