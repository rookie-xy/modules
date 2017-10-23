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
    "github.com/rookie-xy/hubble/plugin"
    "github.com/rookie-xy/hubble/output"
    "github.com/rookie-xy/hubble/pipeline"
    "github.com/rookie-xy/modules/proxy/forward/worker"
    "github.com/rookie-xy/hubble/job"
)

const Name  = "forward"

type forward struct {
    log      log.Log

    queue    pipeline.Queue

    worker  *worker.Worker
    jobs    *job.Jobs
}

var (
    batch     = command.New( module.Flag, "batch",    nil, "This option use to group" )
    client    = command.New( plugin.Flag, "client.elasticsearch",    nil, "This option use to group" )
    sincedb   = command.New( module.Flag, "sincedb",    nil, "This option use to group" )
)

var commands = []command.Item{

    { batch,
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

    { sincedb,
      command.FILE,
      module.Proxy,
      command.SetObject,
      state.Enable,
      0,
      nil },
}

func New(l log.Log) module.Template {
    return &forward{
        log: l,
    }
}

func (r *forward) Init() {
    key := client.GetFlag() + "." + client.GetKey()
    if client, err := factory.Client(key, r.log, client.GetValue()); err != nil {
        fmt.Println("client error ", err)
        return
    } else {
        r.client = client
        register.Forword(key, client)
    }

    return
}

func (f *forward) Main() {
    if f.client == nil || f.pipeline == nil {
        return
    }

    fmt.Println("Start proxy forward module ...")

    for {
        event, err := f.queue.Dequeue(10)
        switch err {

        }

        f.jobs.Start(f.worker)
    }
}

func (r *forward) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Proxy, Name, commands, New)
}

/*
//    pipeline  = command.New( plugin.Flag, "pipeline.stream",  nil, "This option use to group" )
    { pipeline,
      command.FILE,
      module.Proxy,
      command.SetObject,
      state.Enable,
      0,
      nil },


    key := pipeline.GetFlag() + "." + pipeline.GetKey()
    pipeline, err := factory.Pipeline(key, r.log, pipeline.GetValue())
    if err != nil {
        fmt.Println("pipeline error ", err)
        return
    } else {
        r.pipeline = pipeline
    }

    register.Queue(client.GetKey(), pipeline)
*/
