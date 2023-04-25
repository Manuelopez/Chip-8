package main

import (
	"chip-8/chip"
)

const (
    screenWidth  = 64
    screenHeight = 32
)

func main() {
    keyboard := make(map[byte]byte)
    keyboard['1'] = 0;
    keyboard['2'] = 1;
    keyboard['3'] = 2;
    keyboard['4'] = 3;
    keyboard['q'] = 4;
    keyboard['w'] = 5;
    keyboard['e'] = 6;
    keyboard['r'] = 7;
    keyboard['a'] = 8;
    keyboard['s'] = 9;
    keyboard['d'] = 10;
    keyboard['f'] = 11;
    keyboard['z'] = 12;
    keyboard['x'] = 13;
    keyboard['c'] = 14;
    keyboard['v'] = 15;

    c := chip.New(false, keyboard)
    
    //    test display
    // n := "./ibm_logo.ch8"
    //n := "./2-ibm-logo.ch8"

    //    test codes
    //n := "./test_opcode.ch8"
    //n := "./test2.ch8"

    //    falgs
    //n := "./4-flags.ch8"

    //    quirks
    //n := "./5-quirks.ch8"

    //    keyboard
    //n := "./6-keypad.ch8"

    // game

    //n := "./invaders.rom"
    //n := "./connect4.rom"
    n := "./pong.rom"
    

    c.LoadRom(n, 0x200)
    c.Start()
   // c.Test()
}

