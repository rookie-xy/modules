package agents

import (
    "fmt"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/types/value"
    "github.com/rookie-xy/hubble/configure"

  _ "github.com/rookie-xy/modules/agents/file"
)

const Name = module.Agents

var (
    agents = command.New(module.Flag, Name, nil, "agents may be many")
)

var commands = []command.Item{

    { agents,
      command.FILE,
      Name,
      Name,
      nil,
      nil },

}

type Agent struct {
    log.Log

    event     chan int
    children  []module.Template
    done      chan struct{}
}

func New(log log.Log) module.Template {
    new := &Agent{
        Log: log,
        event: make(chan int, 1),
    }

    register.Observer(Name, new)

    return new
}

func (r *Agent) Update(o types.Object) error {
    v := value.New(o)
    if v.GetType() != types.MAP {
        return fmt.Errorf("type is no equeal map")
    }

    if value := v.GetMap(); value != nil {
        val, exist := value[Name]
        if !exist {
            return fmt.Errorf("Not found agents configure")
        }

        agents.SetValue(val)
    }

    r.event <- 1
    return nil
}

func (r *Agent) Init() {
    // 等待配置更新完成的信号
    <-r.event
    fmt.Println("agents init")

    if agents := agents.GetValue(); agents != nil {

        iterator := agents.GetIterator(nil)
        if iterator != nil {
            for {
                if build := configure.Build; build != nil {
                    if err := build(Name, iterator, r.Load); err != nil {
                        fmt.Println("agents init error: ", err)
                        r.Exit(0)
                        return
                    } else {
                        // debug
                        //fmt.Println("agents init not error")
                        break
                    }

                } else {
                    // Debug
                    // fmt.Println("proxy hava not init finish")
                    continue
                }
           }
        }
    }

    return
}

func (r *Agent) Main() {
    fmt.Println("Start agent modules ...")

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
            fmt.Println("agent module exit")
            return
        }
    }
}

func (r *Agent) Exit(code int) {
    defer close(r.done)

    if n := len(r.children); n > 0 {
        for _, child := range r.children {
            child.Exit(code)
        }
    }
}

func (r *Agent) Load(m module.Template) {
    if m != nil {
        r.children = append(r.children, m)
    }
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
