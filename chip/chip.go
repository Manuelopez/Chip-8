package chip

import (
	"chip-8/display"
	"chip-8/keyboard"
	"chip-8/memory"
	"chip-8/register"
	"chip-8/stack"
	"chip-8/timer"
	"chip-8/util"

	//"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Chip struct {
	Memory      *memory.Memory
	Display     *display.Display
	Stack       *stack.Stack
	Keyboard    *keyboard.Keyboard
	PC          uint
	I           [2]*register.Register
	Delay       *timer.Timer
	Sound       *timer.Timer
	RegisterMap map[int]*register.Register

	old         bool
}

func New(old bool, keyboardCgf map[byte]byte) *Chip {
	c := Chip{
		Memory:   memory.New(),
		Display:  display.New(),
		Stack:    stack.New(),
		Keyboard: keyboard.New(keyboardCgf),
		PC:       0,
		I:        [2]*register.Register{register.New(), register.New()},
		Delay:    timer.New(),
		Sound:    timer.New(),
		old:      old,
	}
	c.RegisterMap = make(map[int]*register.Register)
	c.RegisterMap[0] = register.New()
	c.RegisterMap[1] = register.New()
	c.RegisterMap[2] = register.New()
	c.RegisterMap[3] = register.New()
	c.RegisterMap[4] = register.New()
	c.RegisterMap[5] = register.New()
	c.RegisterMap[6] = register.New()
	c.RegisterMap[7] = register.New()
	c.RegisterMap[8] = register.New()
	c.RegisterMap[9] = register.New()
	c.RegisterMap[10] = register.New()
	c.RegisterMap[11] = register.New()
	c.RegisterMap[12] = register.New()
	c.RegisterMap[13] = register.New()
	c.RegisterMap[14] = register.New()
	c.RegisterMap[15] = register.New()

	return &c
}

func (c *Chip) LoadRom(fileName string, startAddr uint16) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	c.PC = uint(startAddr)
	for _, d := range data {

		c.Memory.Write(startAddr, d)
		startAddr++
	}

	return nil
}

func (c *Chip) Start() {
    c.Display.Start()
    go c.KeyPressed()

    firstTimer := true
    ticker := time.NewTicker(time.Microsecond * 1000)
	for range ticker.C{
        if firstTimer == true{
            go c.TimersLoop()
            firstTimer = false
        }

		hbits, lbits := c.Fetch()
		first, second, third, fourth, full := c.Decode(hbits, lbits)
		c.Execute(first, second, third, fourth, full)
	}
}

func (c *Chip) TimersLoop(){
    ticker := time.NewTicker(time.Millisecond * 2)
	for range ticker.C{
        delay := util.BinaryToDecilam8(c.Delay.Get())
        sound := util.BinaryToDecilam8(c.Sound.Get())

        delay--
        sound--

        _, delayBit := util.DecimalToBinary16(uint16(delay))
        _, soundBit := util.DecimalToBinary16(uint16(sound))
        c.Delay.Set(delayBit)
        c.Sound.Set(soundBit)
	}
}

func (c *Chip) Fetch() ([8]bool, [8]bool) {
	first := c.Memory.Read(uint16(c.PC))
	second := c.Memory.Read(uint16(c.PC + 1))
	c.PC += 2

	return first, second
}

func (c *Chip) KeyPressed() {

	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()
    ticker := time.NewTicker(time.Millisecond)

    first := true
    lastChar := rune(0)
	for range ticker.C{
		select {
		case ev := <-events:
			if ev.Ch == 'p' {
				termbox.Close()
				return
			}

            c.Keyboard.SetKey(byte(ev.Ch), true)

            if first || lastChar != ev.Ch{
                first = false
            }

            lastChar = ev.Ch

		default:
            c.Keyboard.SetAllToFalse()
            first = true
            lastChar = 0
		}
	}

}

func (c *Chip) Test() {
	termbox.Init()
	go c.KeyPressed()
	for {
		time.Sleep(time.Millisecond * 40)
	}
}

func (c *Chip) Decode(hbits, lbits [8]bool) (int64, int64, int64, int64, int64) {

	a := make([]bool, 0)
	for i := 0; i < 16; i++ {
		if i < 8 {
			a = append(a, hbits[i])
		} else {
			a = append(a, lbits[i-8])
		}
	}

	hVal := util.BinaryToDecilam8(hbits)
	lVal := util.BinaryToDecilam8(lbits)
	fullVal := util.BinaryToDecilam(a)

	nibble1 := hVal >> 4
	nibble2 := hVal & 0x0F
	nibble3 := lVal >> 4
	nibble4 := lVal & 0x0F

	return nibble1, nibble2, nibble3, nibble4, fullVal
}

func (c *Chip) Execute(first, second, third, fourth, full int64) {

	switch first {
	case 0:
		switch {
		case second == 0 && third == 0xE && fourth == 0:
			c.Display.Screen = [64][32]bool{}
			c.Display.Update()
		case second == 0 && third == 0xE && fourth == 0xE:
			high, low := c.Stack.Pop()
            allbits := make([]bool, 0)
            for i:= 0; i < 16; i++{
                if i < 8{
                    allbits = append(allbits, high[i])
                }else{

                    allbits = append(allbits, low[i-8])
                }
            }
			val := util.BinaryToDecilam(allbits)
            

			c.PC = uint(val)
		}
	case 1:
		c.PC = uint(full) & 0x0FFF
	case 2:
		c.Stack.Push(uint16(c.PC))
		c.PC = uint(full) & 0x0FFF
	case 3:
		val := c.RegisterMap[int(second)].Read()
		decVal := util.BinaryToDecilam8(val)
		if decVal == full&0x00FF {
			c.PC += 2
		}
	case 4:
		val := c.RegisterMap[int(second)].Read()
		decVal := util.BinaryToDecilam8(val)
		if decVal != full&0x00FF {
			c.PC += 2
		}
	case 5:
		val1 := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
		val2 := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
		if val1 == val2 {
			c.PC += 2
		}
	case 6:
		_, val := util.DecimalToBinary16(uint16(full & 0x00FF))
		c.RegisterMap[int(second)].Write(val)
	case 7:
		val1 := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
		val2 := full & 0x00FF
		sum := val1 + val2
		_, reg := util.DecimalToBinary16(uint16(sum))
		c.RegisterMap[int(second)].Write(reg)
	case 8:
		switch fourth {
		case 0:
			regy := c.RegisterMap[int(third)].Read()
			c.RegisterMap[int(second)].Write(regy)
		case 1:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regx | regy
			_, resultB := util.DecimalToBinary16(uint16(result))

			c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
			c.RegisterMap[int(second)].Write(resultB)
		case 2:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regx & regy
			_, resultB := util.DecimalToBinary16(uint16(result))

			c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
			c.RegisterMap[int(second)].Write(resultB)
		case 3:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regx ^ regy
			_, resultB := util.DecimalToBinary16(uint16(result))

			c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
			c.RegisterMap[int(second)].Write(resultB)
		case 4:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regx + regy
			_, resultB := util.DecimalToBinary16(uint16(result))
			c.RegisterMap[int(second)].Write(resultB)
			if result > 255 {
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, true})
			} else {
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
			}
		case 5:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regx - regy
			_, resultB := util.DecimalToBinary16(uint16(result))
			c.RegisterMap[int(second)].Write(resultB)

			if result >= 0 {
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, true})
			} else {
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
			}
		case 6:
			if c.old {
				vyVal := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
				bitShifted := vyVal & 0x1
				shiftedVal := vyVal >> 1

				_, newVx := util.DecimalToBinary16(uint16(shiftedVal))
				c.RegisterMap[int(second)].Write(newVx)

				_, newVf := util.DecimalToBinary16(uint16(bitShifted))
				c.RegisterMap[0xF].Write(newVf)

			} else {
				vxVal := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
				bitShifted := vxVal & 0x1
				shiftedVal := vxVal >> 1
				_, newVx := util.DecimalToBinary16(uint16(shiftedVal))
				c.RegisterMap[int(second)].Write(newVx)

				_, newVf := util.DecimalToBinary16(uint16(bitShifted))
				c.RegisterMap[0xF].Write(newVf)

			}
		case 0xE:
			if c.old {
				vyVal := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
				bitShifted := vyVal >> 7
				shiftedVal := vyVal << 1



				_, newVx := util.DecimalToBinary16(uint16(shiftedVal))
				c.RegisterMap[int(second)].Write(newVx)

				_, newVf := util.DecimalToBinary16(uint16(bitShifted))
				c.RegisterMap[0xF].Write(newVf)


			} else {
				vxVal := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
				bitShifted := vxVal >> 7
				shiftedVal := vxVal << 1
				_, newVx := util.DecimalToBinary16(uint16(shiftedVal))
				c.RegisterMap[int(second)].Write(newVx)

				_, newVf := util.DecimalToBinary16(uint16(bitShifted))
				c.RegisterMap[0xF].Write(newVf)

			}

		case 7:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regy - regx
			_, resultB := util.DecimalToBinary16(uint16(result))
			c.RegisterMap[int(second)].Write(resultB)

			if result >= 0 {
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, true})
			} else {
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
			}
		}

	case 9:
		val1 := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
		val2 := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
		if val1 != val2 {
			c.PC += 2
		}
	case 0xA:
		val := full & 0x0FFF
		hbits, lbits := util.DecimalToBinary16(uint16(val))
		c.I[0].Write(hbits)
		c.I[1].Write(lbits)
	case 0xB:
		if c.old {
			address := full & 0x0FFF
			offset := util.BinaryToDecilam8(c.RegisterMap[0].Read())
			c.PC = uint(address) + uint(offset)
		} else {
			address := full & 0x0FFF
			offset := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			c.PC = uint(address) + uint(offset)
		}

	case 0xC:
		randVal := int64(rand.Intn(256))
		result := randVal & (full & 0x00FF)
		_, bR := util.DecimalToBinary16(uint16(result))
		c.RegisterMap[int(second)].Write(bR)
	case 0xD:
		regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
		regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
		xcoor := regx & 63
		originalX := xcoor
		ycoor := regy & 31
		c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
		hbits := c.I[0].Read()
		lbits := c.I[1].Read()
		ibits := make([]bool, 0)
		for i := 0; i < 16; i++ {
			if i < 8 {
				ibits = append(ibits, hbits[i])
			} else {

				ibits = append(ibits, lbits[i-8])
			}
		}

		for i := 0; i < int(fourth); i++ {
			address := util.BinaryToDecilam(ibits)
			sprite := c.Memory.Read(uint16(address) + uint16(i))

			for _, bit := range sprite {
				if bit && c.Display.Screen[xcoor][ycoor] {
					c.Display.Screen[xcoor][ycoor] = false
					c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, true})
				} else if bit && !c.Display.Screen[xcoor][ycoor] {
					c.Display.Screen[xcoor][ycoor] = bit
				}
				if xcoor >= 63 {
					break
				}
				xcoor++

			}
			xcoor = originalX

			ycoor++

			if ycoor == 32 {
				break
			}

		}
        

        c.Display.Update()
	case 0xE:
		// TODO need to check for keys
        switch{
        case third == 9 && fourth == 0xE:
            //fmt.Println("vx postive", second, " getting the key", c.Keyboard.GetKey(byte(second)))
            key := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
            if c.Keyboard.GetKey(byte(key)){
                c.PC += 2
            }
        case third == 0xA && fourth == 0x1:

            //fmt.Println("vx negative", second, " getting the key", c.Keyboard.GetKey(byte(second)))

            key := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
            if !c.Keyboard.GetKey(byte(key)){
                //fmt.Println("i happend")
                c.PC += 2
            }
        }
        

	case 0xF:
		switch {
		case third == 0 && fourth == 7:
			// TODO timer
            val := c.Delay.Get()
            c.RegisterMap[int(second)].Write(val)
		case third == 1 && fourth == 5:
			// TODO timer

            val := c.RegisterMap[int(second)].Read()
            c.Delay.Set(val)
		case third == 1 && fourth == 8:
			// TODO timer
            val := c.RegisterMap[int(second)].Read()
            c.Delay.Set(val)

		case third == 1 && fourth == 0xE:
			hIbit := c.I[0].Read()
			lIbit := c.I[1].Read()
			iReg := make([]bool, 0)

			for i := 0; i < 16; i++ {
				if i < 8 {
					iReg = append(iReg, hIbit[i])
				} else {
					iReg = append(iReg, lIbit[i-8])
				}

				val := util.BinaryToDecilam(iReg)
				regxVal := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
				sum := val + regxVal
				nhbit, nlbit := util.DecimalToBinary16(uint16(sum))
				c.I[0].Write(nhbit)
				c.I[1].Write(nlbit)
				if !c.old {
					if sum > 0x0FFF {
						c.RegisterMap[0xF].Write([8]bool{true, true, true, true, true, true, true, true})
					}
				}
			}

		case third == 0 && fourth == 0xA:
            // TODO get key pressed

            //fmt.Println("vx pressed", second, "getting the key", c.Keyboard.GetKeyPressed())
            val := c.Keyboard.GetKeyPressed()
            if val == 20{
                c.PC -= 2
            }else{
                _, lbits := util.DecimalToBinary16(uint16(val))
                c.RegisterMap[int(second)].Write(lbits)
            }
		case third == 2 && fourth == 0x9:
			// TODO font
            c.I[1].Write(c.RegisterMap[int(second)].Read())
		case third == 3 && fourth == 3:
			vxVal := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			i2 := vxVal % 10
			vxVal = vxVal / 10
			i1 := vxVal % 10
			vxVal = vxVal / 10
			i0 := vxVal % 10

			hbit, lbit := c.I[0].Read(), c.I[1].Read()
			iBits := make([]bool, 0)
			for i := 0; i < 16; i++ {
				if i < 8 {
					iBits = append(iBits, hbit[i])
				} else {
					iBits = append(iBits, lbit[i-8])
				}
			}

			address := util.BinaryToDecilam(iBits)
			c.Memory.Write(uint16(address), uint8(i0))
			c.Memory.Write(uint16(address+1), uint8(i1))
			c.Memory.Write(uint16(address+2), uint8(i2))

		case third == 5 && fourth == 5:
			hbit, lbit := c.I[0].Read(), c.I[1].Read()
			iBits := make([]bool, 0)
			for i := 0; i < 16; i++ {
				if i < 8 {
					iBits = append(iBits, hbit[i])
				} else {
					iBits = append(iBits, lbit[i-8])
				}
			}

			address := util.BinaryToDecilam(iBits)
			for i := 0; i <= int(second); i++ {
				vx := util.BinaryToDecilam8(c.RegisterMap[i].Read())
				c.Memory.Write(uint16(address), uint8(vx))
				address++ 
			}
			if c.old {
				hnbit, lnbit := util.DecimalToBinary16(uint16(address))
				c.I[0].Write(hnbit)
				c.I[1].Write(lnbit)
			}
		case third == 6 && fourth == 5:
			hbit, lbit := c.I[0].Read(), c.I[1].Read()
			iBits := make([]bool, 0)
			for i := 0; i < 16; i++ {
				if i < 8 {
					iBits = append(iBits, hbit[i])
				} else {
					iBits = append(iBits, lbit[i-8])
				}
			}

			address := util.BinaryToDecilam(iBits)
			for i := 0; i <= int(second); i++ {
				bits := c.Memory.Read(uint16(address))
				c.RegisterMap[i].Write(bits)
				address++
			}
			if c.old {
				hnbit, lnbit := util.DecimalToBinary16(uint16(address))
				c.I[0].Write(hnbit)
				c.I[1].Write(lnbit)
			}

		}

	}

}
