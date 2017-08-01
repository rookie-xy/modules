package channels

import (
    "fmt"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/prototype"
    "github.com/rookie-xy/worker/src/state"
)

const Name = module.Channels

var (
    channels = &command.Meta{ "", Name, nil, "inputs may be many" }
)

var commands = []command.Item{

    { channels,
      command.FILE,
      Name,
      nil,
      0,
      nil },

}

type Channel struct {
    log.Log
    event chan int
    children []module.Template
}

func New(log log.Log) module.Template {
    new := &Channel{
        Log: log,
        event: make(chan int),
    }

    return new
}

func (r *Channel) Update(name string, configure prototype.Object) int {
    if name == "" || configure == nil {
        return state.Error
    }

    if name != Name {
        return state.Declined
    }

    channels.Value = configure
    r.event <-1

    return state.Ok
}

func (r *Channel) Init() {
    fmt.Println("channel init")
    // 等待配置更新完成的信号
    <-r.event

    return
}

func (r *Channel) Main() {
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

func (r *Channel) Exit(code int) {
    /*
    for _, module := range r.children {
        module.Exit()
    }

    //r.cycle.Quit()
    return
    */
}

func (r *Channel) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
