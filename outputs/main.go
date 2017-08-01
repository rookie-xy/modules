package outputs

import (
    "fmt"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/prototype"
    "github.com/rookie-xy/worker/src/state"
)

const Name = module.Outputs

var (
    outputs = &command.Meta{ "", Name, nil, "outputs may be many" }
)

var commands = []command.Item{

    { outputs,
      command.FILE,
      Name,
      nil,
      0,
      nil },

}

type Output struct {
    log.Log
    event chan int
    children []module.Template
}

func New(log log.Log) module.Template {
    new := &Output{
        Log: log,
        event: make(chan int),
    }

    return new
}

func (r *Output) Update(name string, configure prototype.Object) int {
    if name == "" || configure == nil {
        return state.Error
    }

    if name != Name {
        return state.Declined
    }

    outputs.Value = configure
    r.event <-1

    return state.Ok
}

func (r *Output) Init() {
    fmt.Println("output init")
    // 等待配置更新完成的信号
    <-r.event

    return
}

func (r *Output) Main() {
    return
}

func (r *Output) Exit(code int) {
    return
}

func (r *Output) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
