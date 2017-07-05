package inputs

import (
    /*
    "fmt"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/instance"
    "github.com/rookie-xy/worker/src/prototype"
    */
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/prototype"
  _ "github.com/rookie-xy/modules/inputs/src/file"
    "github.com/rookie-xy/worker/src/cycle"
    "github.com/rookie-xy/worker/src/instance"
)

const Name  = "inputs"


type Input struct {
    observers []observer.Observer
    Data prototype.Object
    log.Log
    children []module.Template
}

func New(log log.Log) *Input {
    return &Input{
        Log: log,
    }
}

func (r *Input) Update(configure prototype.Object) {
    if configure != nil {
        inputs.Value <- configure
    }
}

var (
    inputs = command.Meta{ "", Name, nil, "inputs may be many" }
)

var commands = []command.Item{

    { &inputs,
      command.LINE,
      module.GLOBEL,
      nil,
      0,
      nil },

}

func Init(log log.Log) module.Template {
    // 等待配置更新完成的信号

    // TODO load 各个组件
    if v := inputs.Value; v != nil {
        // key为各个模块名字，value为各个模块配置
        for name, configure := range v.(map[string]prototype.Object) {
            // 渲染模块命令
            for key, value := range configure.(map[string]prototype.Object) {
                fmt.Println(key, value)
            }

            if m, ok := module.Pool[name]; ok {
                // TODO 判断作用域
                r.Load(m.Init())
            }
        }

    } else {
        fmt.Println("input value is nil")
    }

    return nil
}

func (r *Input) Main() {
    /*
    // 启动各个组件
    for _, child := range r.children {
        child.Main()
    }
    for ;; {
        //发送消息到channel

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

func (r *Input) Exit() {
    /*
    for _, module := range r.children {
        module.Exit()
    }

    //r.cycle.Quit()
    return
    */
}

func (r *Input) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    instance.Register(Name, Name, commands, Init)
}
