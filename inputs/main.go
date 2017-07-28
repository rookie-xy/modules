package inputs

import (
"fmt"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/prototype"

  _ "github.com/rookie-xy/modules/inputs/src/file"
        "github.com/rookie-xy/worker/src/register"
"github.com/rookie-xy/worker/src/state"
)

const Name = module.Inputs

var (
    inputs = &command.Meta{ "", Name, nil, "inputs may be many" }
)

var commands = []command.Item{

    { inputs,
      command.FILE,
      Name,
      nil,
      0,
      nil },

}

type Input struct {
    log.Log
    event chan int
    children []module.Template
}

func New(log log.Log) module.Template {
    new := &Input{
        Log: log,
        event: make(chan int),
    }

    register.Observer(Name, new)

    return new
}

func (r *Input) Update(name string, configure prototype.Object) int {
    if name == "" || configure == nil {
        return state.Error
    }

    if name != Name {
        return state.Declined
    }

    inputs.Value = configure
    r.event <-1

    return state.Ok
}

func (r *Input) Init() {
    fmt.Println("input init")
    // 等待配置更新完成的信号
    <-r.event

    // TODO load 各个组件
    /*
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
    */

    return
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

func (r *Input) Exit(code int) {
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
    register.Module(module.Worker, Name, commands, New)
}

/*
func (r *Input) Build(log log.Log) *Input {
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

func (r *Input) Load(m module.Template) {
    r.children = append(r.children, m)
}
*/
