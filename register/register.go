package register

type Register struct{
   V [8]bool
}

func New() *Register{
    return &Register{}
}

func (r *Register) Write(value [8]bool){
    r.V = value
}

func (r Register) Read() [8]bool{
    return r.V
}

