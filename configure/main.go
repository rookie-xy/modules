package configure

import (
    "fmt"
    "unsafe"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/instance"

  _ "github.com/rookie-xy/modules/configure/src/file"
  _ "github.com/rookie-xy/modules/configure/src/zookeeper"
)

const Name  = "configure"

var (
    config = &command.Meta{ "-c", "config", "file", "This configure file path" }
)

var commands = []command.Item{

    { config,
      command.LINE,
      module.GLOBEL,
      command.SetObject,
      unsafe.Offsetof(config.Value),
      nil },

}

type Configure struct {
    log.Log
    children []module.ModuleTemplate
}

func New(log log.Log) *Configure {
    return &Configure{
        Log: log,
    }
}

func (r *Configure) Init() {
    // 根据指令加载所需模块
    name := Name
    if v := config.Value; v != nil {
        name = v.(string)
    } else {
        return
    }

    fmt.Println(name)

    if m, ok := module.Pool[name]; ok {
        r.Load(m)
    }

    return
}

func (r *Configure) Main() {
    // 启动各个子模块组件
    for _, child := range r.children {
        /*
        if child.Init() == Error {
            return
        }
        */
        fmt.Println("qqqqqqqqqqqqqqqq")

        child.Init()
        child.Main()
    }
/*
    for ;; {
        // TODO 监听外部启停指令
        select {
        case <-r.cycle.Stop():
        //child.Exit()
            cycle.Stop()

        default:

        }
    }
    */

    return
}

func (r *Configure) Exit() {
    for _, module := range r.children {
        module.Exit()
    }

    //r.cycle.Quit()
    return
}

func (r *Configure) Load(m module.ModuleTemplate) {
    r.children = append(r.children, m)
}

func init() {
    instance.Register(Name, commands, nil)
}
