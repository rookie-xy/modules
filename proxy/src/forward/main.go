package forward

import (
    "fmt"

    "github.com/rookie-xy/hubble/src/command"
    "github.com/rookie-xy/hubble/src/module"
    "github.com/rookie-xy/hubble/src/log"
    "github.com/rookie-xy/hubble/src/register"
    "github.com/rookie-xy/hubble/src/factory"
    "github.com/rookie-xy/hubble/src/state"
  cli "github.com/rookie-xy/hubble/src/client"
        "github.com/rookie-xy/hubble/src/plugin"
  pipe "github.com/rookie-xy/hubble/src/pipeline"
)

const Name  = "forward"

type forward struct {
    log.Log
    client  cli.Client
    pipeline pipe.Pipeline
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
    key := pipeline.GetFlag() + "." + pipeline.GetKey()
    if pipeline, err := factory.Pipeline(key, r.Log, pipeline); err != nil {
        fmt.Println("pipeline error", err)
        return
    } else {
        r.pipeline = pipeline
    }

    key = client.GetFlag() + "." + client.GetKey()
    if client, err := factory.Client(key, r.Log, client); err != nil {
        fmt.Println("client error", err)
        return
    } else {
        r.client = client
    }

    return
}

func (r *forward) Main() {
    if r.client == nil || r.pipeline == nil {
        return
    }

    fmt.Println("Start proxy forward module ...")

    for {
        event, status := r.pipeline.Pull(10)

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
