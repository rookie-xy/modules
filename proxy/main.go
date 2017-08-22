package outputs

import (
    "fmt"

    "github.com/rookie-xy/hubble/src/command"
    "github.com/rookie-xy/hubble/src/module"
    "github.com/rookie-xy/hubble/src/log"
    "github.com/rookie-xy/hubble/src/register"
    "github.com/rookie-xy/hubble/src/state"
    "github.com/rookie-xy/hubble/src/types"
	"github.com/rookie-xy/hubble/src/types/value"

  _ "github.com/rookie-xy/modules/proxy/src/forward"
  _ "github.com/rookie-xy/modules/proxy/src/reverse"

    "github.com/rookie-xy/hubble/src/configure"
)

const Name = module.Proxy

var (
    proxy = command.New("", Name, nil, "outputs may be many")
)

var commands = []command.Item{

    { proxy,
      command.FILE,
      Name,
      nil,
      state.Enable,
      0,
      nil },

}

type Proxy struct {
    log.Log
    event    chan int
    children []module.Template
}

func New(log log.Log) module.Template {
    new := &Proxy{
        Log: log,
        event: make(chan int, 1),
    }

    register.Observer(Name, new)

    return new
}

func (r *Proxy) Init() {
    <-r.event
    fmt.Println("proxy init")

    build := func(name string, i types.Iterator, load module.Load) int {
        for iterator := i; iterator.Has(); iterator.Next() {
            iterm := iterator.Iterm()

            if v := iterm.Value; v != nil {
                value := value.New(v)
                it := value.GetIterator(nil)
                if it == nil {
                    continue
                }

                for iterator := it; iterator.Has(); iterator.Next() {
                    if iterm := iterator.Iterm(); iterm != nil {
                        key := iterm.Key.GetString()

                        if status := command.File(name, key, iterm.Value); status != state.Ok {
                            fmt.Println("command file error", status)
                            return state.Error
                        }
                    }
                }
            }

            namespace := name + "." + iterm.Key.GetString()
            module := module.Setup(namespace, r.Log)
            if module != nil {
                module.Init()

            } else {
                fmt.Printf("[%s] module setup error", name)
                return state.Error
            }

            load(module)
        }

        return state.Ok
    }

    if proxy := proxy.GetValue(); proxy != nil {
        iterator := proxy.GetIterator(nil)
        if iterator != nil {
            if build(Name, iterator, r.Load) == state.Error {
                return
            }

            configure.Build = build
        }
    }

    return
}

func (r *Proxy) Main() {
    fmt.Println("Start proxy modules ...")
    if len(r.children) < 1 {
        return
    }

    for _, child := range r.children {
        child.Main()
    }

    for ;; {
        select {}
    }
}

func (r *Proxy) Exit(code int) {
    return
}

func (r *Proxy) Update(o types.Object) int {
    v := value.New(o)
    if v.GetType() != types.Map {
        return state.Error
    }

    if value := v.GetMap(); value != nil {
        val, exist := value[Name]
        if !exist {
            fmt.Println("Not found proxy configure")
            return state.Error
        }

        proxy.SetValue(val)
    }

    r.event <- 1
    return state.Ok
}

func (r *Proxy) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
