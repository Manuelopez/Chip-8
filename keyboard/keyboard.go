package keyboard

type Keyboard struct{
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
    
func (k *Keyboard) GetKey(key byte) bool{
    val := k.CFG[key]
    return k.Keys[val]
}

func (k *Keyboard) SetKey(key byte, val bool){
    address := k.CFG[key]
    k.Keys[address] = val
}


