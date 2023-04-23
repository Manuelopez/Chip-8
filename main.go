package main

import (
	"chip-8/chip"
)

const (
    screenWidth  = 64
    screenHeight = 32
)

func main() {
    c := chip.New(false)
    
   // n := "./ibm_logo.ch8"
    n := "./2-ibm-logo.ch8"

    //n := "./test_opcode.ch8"
    c.LoadRom(n, 0x200)
    c.Start()



}

