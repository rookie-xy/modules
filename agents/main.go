package agents

import (
    "fmt"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/state"
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/types/value"
    "github.com/rookie-xy/hubble/configure"

  _ "github.com/rookie-xy/modules/agents/log/file"
  _ "github.com/rookie-xy/modules/agents/log/stdin"
  _ "github.com/rookie-xy/modules/agents/fpm"
  _ "github.com/rookie-xy/modules/agents/proc"
)

const Name = module.Agents

var (
    agents = command.New(module.Flag, Name, nil, "inputs may be many")
)

var commands = []command.Item{

    { agents,
      command.FILE,
      Name,
      nil,
      state.Enable,
      0,
      nil },

}

type Agent struct {
    log.Log
    event chan int
    children []module.Template
}

func New(log log.Log) module.Template {
    new := &Agent{
        Log: log,
        event: make(chan int, 1),
    }

    register.Observer(Name, new)

    return new
}

func (r *Agent) Update(o types.Object) int {
    v := value.New(o)
    if v.GetType() != types.Map {
        return state.Error
    }

    if value := v.GetMap(); value != nil {
        val, exist := value[Name]
        if !exist {
            fmt.Println("Not found agents configure")
            return state.Error
        }

        agents.SetValue(val)
    }

    r.event <- 1
    return state.Ok
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
                    if build(Name, iterator, r.Load) == state.Error {
                        fmt.Println("agents init error")
                        return
                    } else {
                        fmt.Println("agents init not error")
                        break
                    }

                } else {
                    fmt.Println("proxy hava not init finish")
                    continue
                }
            }
        }
    }

    return
}

func (r *Agent) Main() {
    fmt.Println("Start agent modules ...")
    if len(r.children) < 1 {
        return
    }

    for _, child := range r.children {
        go child.Main()
    }

    for ;; {
        select {}
    }
}

func (r *Agent) Exit(code int) {
    /*
    for _, module := range r.children {
        module.Exit()
    }

    //r.cycle.Quit()
    return
    */
}

func (r *Agent) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
