package timer

import (
	"chip-8/register"
	"sync"
)

type Timer struct{
    mu sync.Mutex
    Reg register.Register
}

func New() *Timer{
    t := Timer{}

    t.Reg.Write([8]bool{true, true, true, true, true, true, true, true})

    return &t
}

func (t *Timer) Set(val [8]bool){
    t.mu.Lock()
    defer t.mu.Unlock()

    t.Reg.Write(val)
}

func (t *Timer) Get() [8]bool{

    t.mu.Lock()
    defer t.mu.Unlock()

    return t.Reg.Read()
}
