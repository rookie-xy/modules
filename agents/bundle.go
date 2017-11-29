package agents

import (
    "sync"

    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/configure"

    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/types/value"

    "github.com/rookie-xy/hubble/log"
  . "github.com/rookie-xy/hubble/log/level"

  _ "github.com/rookie-xy/modules/agents/file"
  //  "github.com/rookie-xy/hubble/adapter"
)

var agents = command.New(module.Flag, Name, nil, "agents may be many")

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
    level     Level

    wg        sync.WaitGroup
    event     chan int
    children  []module.Template
    done      chan struct{}
}

func New(log log.Log) module.Template {
    new := &Agent{
        Log: log,
        //level: adapter.ToLevelLog(log).Level(),
        event: make(chan int, 1),
    }

    register.Observer(Name, new)
    return new
}

func (r *Agent) Update(o types.Object) error {
    v := value.New(o)
    if v.GetType() != types.MAP {
        return ErrType
    }

    if value := v.GetMap(); value != nil {
        val, exist := value[Name]
        if !exist {
            return ErrConfigure
        }

        agents.SetValue(val)
    }

    r.log(DEBUG,"Update agents configure successful")
    r.event <- 1

    return nil
}

func (r *Agent) Init() {
    r.log(DEBUG, "Waiting for a signal to configure the update")

    <-r.event
    r.log(DEBUG,"Initialization components for agent")

    r.done = make(chan struct{})
    r.children = []module.Template{}

    if agents := agents.GetValue(); agents != nil {

        iterator := agents.GetIterator(nil)
        if iterator != nil {
            for {
                if build := configure.Build; build != nil {
                    if err := build(Name, iterator, r.Load); err != nil {
                        r.log(ERROR,"Agents init error: %s\n", err)
                        r.Exit(0)

                        return

                    } else {
                        r.log(INFO,"Agents configure builder finish")
                        break
                    }

                } else {
                    r.log(WARN,"Proxy have not init finish")
                    continue
                }
           }
        }
    }

    r.log(DEBUG,"Agent all component initialization completed")
}

func (r *Agent) Main() {
    r.log(DEBUG,"Run components for agent")

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

            continue
        }

        r.log(WARN, "Agent child component is nil [main stage]")
    }

    r.log(DEBUG, "Agent all components have started running")
}

func (r *Agent) Exit(code int) {
    defer func() {
        <-r.done
        r.log(DEBUG,"Agent all components have exit")
    }()

    r.log(DEBUG,"Exit components for agent")

    if n := len(r.children); n > 0 {
        for _, child := range r.children {
        	if child != nil {
                child.Exit(code)
                continue
            }

            r.log(WARN, "Agent child component is nil [exit stage]")
        }
    }
}

func (r *Agent) Load(m module.Template) {
    if m != nil {
        r.children = append(r.children, m)
    }
}

func (r *Agent) log(l Level, f string, args ...interface{}) {
    log.Print(r.Log, r.level, l, f, args...)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
