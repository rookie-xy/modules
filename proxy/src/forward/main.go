package forward

import (
    "fmt"

    "github.com/rookie-xy/hubble/src/command"
    "github.com/rookie-xy/hubble/src/module"
    "github.com/rookie-xy/hubble/src/log"
    "github.com/rookie-xy/hubble/src/register"
    "github.com/rookie-xy/hubble/src/state"
//    "github.com/rookie-xy/hubble/src/factory"
//    "github.com/rookie-xy/hubble/src/types"
 cli "github.com/rookie-xy/hubble/src/client"
        "github.com/rookie-xy/hubble/src/plugin"
)

const Name  = "forward"

type forward struct {
    log.Log
    clients []cli.Client
}

var (
    pipeline  = command.New( plugin.Flag, "pipeline.stream",  nil, "This option use to group" )
    client    = command.New( plugin.Flag, "client.kafka",    nil, "This option use to group" )
)

var commands = []command.Item{

    { pipeline,
      command.FILE,
      module.Proxy,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { client,
      command.FILE,
      module.Proxy,
      command.SetObject,
      state.Enable,
      0,
      nil },

}

func New(log log.Log) module.Template {
    return &forward{
        Log: log,
    }
}

func (r *forward) Init() {

    fmt.Println("hhhhhhhhhhhhhhhhhhhhhhhhhhhhh", pipeline.GetKey(), pipeline.GetMap())
    fmt.Println("iiiiiiiiiiiiiiiiiiiiiiiiiiiii", client.GetKey(), client.GetMap())

    fmt.Println(Name + " init")
    /*
    subscribes := subscribe.Value.([]interface{})
    for index, element := range subscribes {
        fmt.Println("subscribessssssssssssss ", index, element)
    }
    */
/*
    clients := client.Value.(map[interface{}]interface{})
    for key, value := range clients {
        if client, err := factory.Client(value.(string)); err == nil {
            r.clients = append(r.clients, client)
            //fmt.Println("factory client error")
        }

        fmt.Println("clientssssssssssssss ", key, value)
    }
    */

    return
}

func (r *forward) Main() {

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

func (r *forward) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Proxy, Name, commands, New)
}
