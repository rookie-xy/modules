package configure

import "github.com/rookie-xy/hubble/observer"

func Init(inits []observer.Observer) {
    for _, init := range inits {
        init.Reinit()
    }
}

func Main(mains []observer.Observer) {
    for _, main := range mains {
        main.Remain()
    }
}
