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
    "sync"
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

    wg        sync.WaitGroup
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
    fmt.Println("Initialization components for agent")
    r.done = make(chan struct{})
    r.children = []module.Template{}

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

    // debug
    fmt.Println("Agent all component initialization completed")
}

func (r *Agent) Main() {
    fmt.Println("Run components for agent")

    defer func() {
        r.wg.Wait()
        close(r.done)
    }()

    r.wg.Add(len(r.children))

    for _, child := range r.children {
        if child != nil {
        	go func(main func()) {
        	    defer r.wg.Done()

        	    main()

            }(child.Main)
        }
    }

    //debug
    fmt.Println("Agent all components have started running")
}

func (r *Agent) Exit(code int) {
    defer func() {
        <-r.done
        fmt.Println("Agent all components have exit")
    }()

    // debug
    fmt.Println("Exit components for agent")

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
