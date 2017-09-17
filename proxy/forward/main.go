package forward

import (
    "fmt"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/state"
    "github.com/rookie-xy/hubble/proxy"
    queue "github.com/rookie-xy/hubble/pipeline"
    "github.com/rookie-xy/hubble/plugin"
)

const Name  = "forward"

type forward struct {
    log.Log
    client   proxy.Forward
    pipeline queue.Queue
}

var (
    pipeline  = command.New( plugin.Flag, "pipeline.stream",  nil, "This option use to group" )
    client    = command.New( plugin.Flag, "client.elasticsearch",    nil, "This option use to group" )
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
    key := pipeline.GetFlag() + "." + pipeline.GetKey()
    pipeline, err := factory.Pipeline(key, r.Log, pipeline.GetValue())
    if err != nil {
        fmt.Println("pipeline error", err)
        return
    } else {
        r.pipeline = pipeline
    }

    register.Queue(client.GetKey(), pipeline)

    key = client.GetFlag() + "." + client.GetKey()
    if client, err := factory.Client(key, r.Log, client.GetValue()); err != nil {
        fmt.Println("client error", err)
        return
    } else {
        r.client = client
        register.Forword(key, client)
    }

    return
}

func (r *forward) Main() {
    if r.client == nil || r.pipeline == nil {
        return
    }

    fmt.Println("Start proxy forward module ...")

    for {
        event, status := r.pipeline.Dequeue(10)

        switch status {

        case state.Ignore:
            continue

        case state.Busy:
            //TODO sleep
        }

        r.client.Sender(event)
    }
}

func (r *forward) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Proxy, Name, commands, New)
}
