package sender

import (
    "fmt"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/state"
    "github.com/rookie-xy/worker/src/factory"
//    "github.com/rookie-xy/worker/src/types"
 cli "github.com/rookie-xy/worker/src/client"
)

const Name  = "sender"

type sender struct {
    log.Log
    clients []cli.Client
}

var (
    subscribe = command.Metas( "", "subscribe", nil, "This option use to group" )
    client    = command.Metas( "", "client",    nil, "This option use to group" )
)

var commands = []command.Item{

    { subscribe,
      command.FILE,
      module.Outputs,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { client,
      command.FILE,
      module.Outputs,
      command.SetObject,
      state.Enable,
      0,
      nil },

}

func New(log log.Log) module.Template {
    return &sender{
        Log: log,
    }
}

func (r *sender) Init() {

    fmt.Println(Name + " init")

    subscribes := subscribe.Value.([]interface{})
    for index, element := range subscribes {
        fmt.Println("subscribessssssssssssss ", index, element)
    }

    clients := client.Value.(map[interface{}]interface{})
    for key, value := range clients {
        if client, err := factory.Client(value); err == nil {
            r.clients = append(r.clients, client)
            //fmt.Println("factory client error")
        }

        fmt.Println("clientssssssssssssss ", key, value)
    }

    return
}

func (r *sender) Main() {

    for {
        select {
/*
        case:

            for _, client := range r.clients {
                client.Sender()
            }
            */
        }
    }
}

func (r *sender) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Outputs, Name, commands, New)
}
