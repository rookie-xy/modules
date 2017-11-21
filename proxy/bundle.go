package proxy

import (
    "fmt"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/configure"
    "github.com/rookie-xy/hubble/types/value"

  _ "github.com/rookie-xy/modules/proxy/forward"
  _ "github.com/rookie-xy/modules/proxy/sinceDB"
)

const Name = module.Proxy

var (
    proxy = command.New("", Name, nil, "outputs may be many")
)

var commands = []command.Item{

    { proxy,
      command.FILE,
      Name,
      Name,
      nil,
      nil },

}

type Proxy struct {
    log.Log
    event     chan int
    children  []module.Template
    done      chan struct{}
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

    build := func(scope string, i types.Iterator, load module.Load) error {
        for iterator := i; iterator.Has(); iterator.Next() {
            iterm := iterator.Iterm()
            name := iterm.Key.GetString()


            if v := iterm.Value; v != nil {
                value := value.New(v)
                it := value.GetIterator(nil)
                if it == nil {
                    continue
                }

                for iterator := it; iterator.Has(); iterator.Next() {
                    if iterm := iterator.Iterm(); iterm != nil {
                        key := iterm.Key.GetString()
                        if err := command.File(scope, name, key, iterm.Value); err != nil {
                            return fmt.Errorf("command file error ", err)
                        }
                    }
                }
            }

            namespace := scope + "." + name
            module := module.Setup(namespace, r.Log)
            if module != nil {
                module.Init()

            } else {
                return fmt.Errorf("[%s] module setup error\n", name)
            }

            load(module)
        }

        return nil
    }

    if proxy := proxy.GetValue(); proxy != nil {
        iterator := proxy.GetIterator(nil)
        if iterator != nil {
            if err := build(Name, iterator, r.Load); err != nil {
                fmt.Println("proxy build error ", err)
                return
            }

            configure.Build = build
        }
    }

    return
}

func (r *Proxy) Main() {
    fmt.Println("Start proxy modules ...")

    if len(r.children) > 0 {
        for i, child := range r.children {
            if child != nil {
                go child.Main()
            } else {
                fmt.Println("error")
                if i > 0 {
                    r.Exit(0)
                }
                return
            }
        }
    } else {
        fmt.Println("ERROR")
        return
    }

    for {
        select {
        case <-r.done:
            fmt.Println("proxy module exit")
            return
        }
    }
}

func (r *Proxy) Exit(code int) {
    defer close(r.done)

    if n := len(r.children); n > 0 {
        for _, child := range r.children {
            if child != nil {
                child.Exit(code)
            }
        }
    }
}

func (r *Proxy) Update(o types.Object) error {
    v := value.New(o)
    if v.GetType() != types.MAP {
        return fmt.Errorf("TYPE is not equal map")
    }

    if value := v.GetMap(); value != nil {
        val, exist := value[Name]
        if !exist {
            return fmt.Errorf("Not found proxy configure")
        }

        proxy.SetValue(val)
    }

    r.event <- 1
    return nil
}

func (r *Proxy) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
