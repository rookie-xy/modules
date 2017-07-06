package configure

import (
    /*
    "fmt"
    "unsafe"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/instance"
    "github.com/rookie-xy/worker/src/configure"
    */

  _ "github.com/rookie-xy/modules/configure/src/file"
  _ "github.com/rookie-xy/modules/configure/src/zookeeper"
    "github.com/rookie-xy/worker/src/cycle"
    "github.com/rookie-xy/worker/src/instance"
    "github.com/rookie-xy/worker/src/module"
    "github.com/elastic/beats/filebeat/prospector/log"
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
    *configure.Configure
    children []module.Template
}

func New(log log.Log) *Configure {
    return &Configure{
        Log: log,
    }
}

func Init(log log.Log) module.Template {
    // 根据指令加载所需模块
    config := New(log)

    name := Name
    if v := config.Value; v != nil {
        name = v.(string)
    } else {
        return
    }

    if m, ok := module.Pool[name]; ok {
        //判断作用域
        config.Load(m(log))
    }

    return config
}

func (r *Configure) Main() {
    // 启动各个子模块组件
/*
    for _, child := range r.children {
        if child.Init() == Error {
            return
        }
        fmt.Println("qqqqqqqqqqqqqqqq")

        child.Init()
        child.Main()
    }

    r.Notify()

    // 渲染所有配置指令
    for ;; {
        // TODO 解析配置，通知加载三大模块
        // TODO 监听外部启停指令
        select {
        case <-r.cycle.Stop():
        //child.Exit()
            cycle.Stop()

        default:

        }
    }

    return
    */
}

func (r *Configure) Exit() {
    for _, module := range r.children {
        module.Exit()
    }

    //r.cycle.Quit()
    return
}

func (r *Configure) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    instance.Register(Name, Name, commands, Init)
}
