package timer

import "chip-8/register"

type Timer struct{
    Reg register.Register
}

func New() *Timer{
    t := Timer{}

    t.Reg.Write([8]bool{true, true, true, true, true, true, true, true})

    return &t
}
