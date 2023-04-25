package keyboard

import "sync"


type Keyboard struct{
    mu sync.Mutex
    CFG map[byte]byte
    Keys map[byte]bool
}

func New(cfg map[byte]byte) *Keyboard{
    address := make(map[byte]bool)
    for _, a := range cfg{
        if a <0 || a > 15{
            panic("address must be between 0 and 15");
        }
        if _, ok := address[a]; ok{
            panic("two keys with the same value");
        }
        address[a] = false
    }
    k := Keyboard{
        CFG: cfg,
        Keys: address,
    }

    return &k

}
    
func (k *Keyboard) GetKey(key byte) (bool){
    k.mu.Lock()
    defer k.mu.Unlock()
    return k.Keys[key]
}

func (k *Keyboard) SetKey(key byte, val bool){

    k.mu.Lock()
    defer k.mu.Unlock()
    address := k.CFG[key]
    k.Keys[address] = val
}

func (k *Keyboard) SetAllToFalse(){
    
    k.mu.Lock()
    defer k.mu.Unlock()
    for keys, _ := range k.Keys{
        k.Keys[keys] = false
    }
}

func (k *Keyboard) GetKeyPressed() (byte){

    k.mu.Lock()
    defer k.mu.Unlock()
    for keys, _ := range k.Keys{
        if k.Keys[keys]{
            return keys
        }
    }

    return 20
}
