package chip

import (
	"chip-8/display"
	"chip-8/memory"
	"chip-8/register"
	"chip-8/stack"
	"chip-8/timer"
	"chip-8/util"

)

type Chip struct {
	Memory      *memory.Memory
	Display     *display.Display
	Stack       *stack.Stack
	PC          uint
	I           [2]*register.Register
	Delay       *timer.Timer
	Sound       *timer.Timer
	RegisterMap map[int]*register.Register
    old bool
}

func New(old bool) *Chip {
	c := Chip{
		Memory:  memory.New(),
		Display: display.New(),
		Stack:   stack.New(),
		PC:      0,
		I:       [2]*register.Register{register.New(), register.New()},
		Delay:   timer.New(),
		Sound:   timer.New(),
        old: old,
	}
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

func (c *Chip) Start() {
	for {
		hbits, lbits := c.Fetch()
		c.Decode(hbits, lbits)
	}
}

func (c *Chip) Fetch() ([8]bool, [8]bool) {
	first := c.Memory.Read(uint16(c.PC))
	second := c.Memory.Read(uint16(c.PC))
	c.PC += 2

	return first, second
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
			counter := c.Stack.Pop()
			val := util.BinaryToDecilam8(counter)
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
			c.RegisterMap[int(second)].Write(resultB)
		case 2:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regx & regy
			_, resultB := util.DecimalToBinary16(uint16(result))
			c.RegisterMap[int(second)].Write(resultB)
		case 3:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regx ^ regy
			_, resultB := util.DecimalToBinary16(uint16(result))
			c.RegisterMap[int(second)].Write(resultB)
		case 4:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regx + regy
			if result > 255 {
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, true})
			}else{
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
            }
			_, resultB := util.DecimalToBinary16(uint16(result))
			c.RegisterMap[int(second)].Write(resultB)
        case 5:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regx - regy
			if result >=0 {
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, true})
			}else{
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
            }
			_, resultB := util.DecimalToBinary16(uint16(result))
			c.RegisterMap[int(second)].Write(resultB)
        case 6:

            

        case 7:
			regx := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
			regy := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
			result := regy - regx  
			if result >=0 {
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, true})
			}else{
				c.RegisterMap[0xF].Write([8]bool{false, false, false, false, false, false, false, false})
            }
			_, resultB := util.DecimalToBinary16(uint16(result))
			c.RegisterMap[int(second)].Write(resultB)
		}

	case 9:
		val1 := util.BinaryToDecilam8(c.RegisterMap[int(second)].Read())
		val2 := util.BinaryToDecilam8(c.RegisterMap[int(third)].Read())
		if val1 != val2 {
			c.PC += 2
		}
	}

}
