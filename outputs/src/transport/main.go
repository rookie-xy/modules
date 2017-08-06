package memory

import (
    "fmt"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/state"

)

const Name  = "transport"

type transportOutput struct{
    log.Log
}

var (
    subscribe = command.Metas( "", "subscribe", nil, "This option use to group" )
    sender    = command.Metas( "", "sender",    nil, "This option use to group" )
)

var commands = []command.Item{

    { subscribe,
      command.FILE,
      module.Outputs,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { sender,
      command.FILE,
      module.Outputs,
      command.SetObject,
      state.Enable,
      0,
      nil },

}

func New(log log.Log) module.Template {
    return &transportOutput{
        Log: log,
    }
}

func (r *transportOutput) Init() {
    fmt.Println("transporttttttttttttttttttt", subscribe.Value, sender.Value)

    return
}

func (r *transportOutput) Main() {
    for {
        select {

        }

    }
}

func (r *transportOutput) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Outputs, Name, commands, New)
}
