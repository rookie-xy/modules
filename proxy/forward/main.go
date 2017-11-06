package forward

import (
    "fmt"
    "strings"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/plugin"
    "github.com/rookie-xy/hubble/pipeline"
    "github.com/rookie-xy/hubble/job"

    "github.com/rookie-xy/modules/proxy/forward/worker"
)

const Name  = "forward"

type forward struct {
    log      log.Log
    queue    pipeline.Queue
    jobs    *job.Jobs

    client  *command.Command
    sinceDB *command.Command
}

var (
    batch     = command.New( module.Flag, "batch",    nil, "This option use to group" )
    client    = command.New( plugin.Flag, "client.kafka",    nil, "This option use to group" )
    sinceDB   = command.New( plugin.Flag, "client.sinceDB",    nil, "This option use to group" )
    pipe      = command.New( plugin.Flag, "pipeline.stream",  nil, "This option use to group" )
)

var commands = []command.Item{

    { batch,
      command.FILE,
      module.Proxy,
      Name,
      command.SetObject,
      nil },

    { pipe,
      command.FILE,
      module.Proxy,
      Name,
      command.SetObject,
      nil },

    { client,
      command.FILE,
      module.Proxy,
      Name,
      command.SetObject,
      nil },

    { sinceDB,
      command.FILE,
      module.Proxy,
      Name,
      command.SetObject,
      nil },
}

func New(log log.Log) module.Template {
    return &forward{
        log:    log,
        jobs:   job.New(log),
    }
}

func (f *forward) Init() {
    key := pipe.GetFlag() + "." + pipe.GetKey()
    pipeline, err := factory.Pipeline(key, f.log, pipe.GetValue())
    if err != nil {
        fmt.Println("pipeline error ", err)
        return
    } else {
        f.queue = pipeline
    }

    name := client.GetKey()
    name = name[strings.LastIndex(name, ".") + 1:]

    register.Queue(name, pipeline)

    f.client = client
    f.sinceDB = sinceDB

    return
}

func (f *forward) Main() {
    fmt.Println("Start proxy forward module ...")

    for {
        event, err := f.queue.Dequeue(10)
        switch err {

        }
        fmt.Println("forwardddddddddddddddddddddddddddddddddddddddd")

        worker := worker.New(f.log)
        if err := worker.Init(f.client, f.sinceDB, event); err != nil {
            fmt.Println(err)
            return
        }

        f.jobs.Start(worker)
    }
}

func (f *forward) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Proxy, Name, commands, New)
}
