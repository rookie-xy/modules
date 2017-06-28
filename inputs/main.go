package inputs

import (
    "github.com/rookie-xy/worker/src/command"
//    "github.com/rookie-xy/worker/src/configure"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/prototype"
    "github.com/rookie-xy/worker/src/cycle"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
)

const Name  = "inputs"

type Input struct {
    children  map[module.Module]prototype.Object
    //children []module.Module
    cycle    cycle.Cycle
    log.Log
}

func New(log log.Log) *Input {
    return &Input{}
}

var (
    inputs = command.Meta{ "", Name, nil, "inputs may be many" }
)

var inputCommands = []command.Item{

    { &inputs,
      command.LINE,
      module.GLOBEL,
      nil,
      0,
      nil },
}

func (r *Input) Init() {
    // 渲染inputs组件命令, 需要原生配置支持
    for key, value := range r.children {
        /*
        reg := register.GetInstance().Module(key, nil, nil)
        r.Load(reg.GetModule(), nil)
        */

        /*
        if child := module.Template(key); child != nil {
            r.Load(child, value)

            if commands := command.Commands(child); commands != nil {
                for cmd := range commands {
                    cmd.SetFunc(cmd.Offset, value[cmd.Name.Key])
                }
            }
        }
        */
    }

    return
}

func (r *Input) Main() {
    cycle := cycle.New()

    // 启动各个组件
    for _, child := range r.children {
        /*
        if child.Init() == Error {
            return
        }
        */

        child.Init()

        //cycle.Start(child.Main(), nil)
    }

    for ;; {
        //发送消息到channel

        select {
/*
        case <-r.cycle.Stop():
        //child.Exit()
            cycle.Stop()

        default:
        */

        }
    }

    return
}

func (r *Input) Exit() {
    for _, module := range r.children {
        module.Exit()
    }

    //r.cycle.Quit()
    return
}

func (r *Input) Load(key module.Module, value map[prototype.Object]prototype.Object) {
    r.children[key] = value
    //r.children = append(r.children, child)
}
