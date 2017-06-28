package inputs

import (
    "hubble/src/command"
    "hubble/src/configure"
    "hubble/src/module"
    "hubble/src/prototype"
)

const Name  = "inputs"

type Input struct {
    options  map[string]prototype.Object
    children []module.Module
    cycle    cycle.Cycle
}

func New() *Input {
    return &Input{}
}

var (
    inputs = command.Name{ "", Name, nil, "inputs may be many" }
)

var inputCommands = []command.Command{

    { inputs,
      configure.COMMAND,
      nil,
      nil,
      nil },
}

func (r *Input) Init() {
    // 渲染inputs组件命令, 需要原生配置支持
    for key, value := range r.options {
        if child := module.Template(key); child != nil {
            r.Load(child, value)

            if commands := command.Commands(child); commands != nil {
                for cmd := range commands {
                    cmd.SetFunc(cmd.Offset, value[cmd.Name.Key])
                }
            }
        }
    }

    return
}

func (r *Input) Main() {
    cycle := cycle.New()

    // 启动各个组件
    for child := range r.children {
        if child.Init() == Error {
            return
        }

        cycle.Start(child.Main(), nil)
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
}

func (r *Input) Exit() {
    for module := range r.children {
        module.Exit()
    }

    //r.cycle.Quit()
    return
}

func (r *Input) Load(child module.Module, value map[string]prototype.Object) {
    r.options[child] = value
    r.children = append(r.children, child)
}
