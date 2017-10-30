package forward

import (
    "fmt"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/register"
 //   "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/state"
//    "github.com/rookie-xy/hubble/proxy"
    "github.com/rookie-xy/hubble/plugin"
//    "github.com/rookie-xy/hubble/output"
    "github.com/rookie-xy/hubble/pipeline"
    "github.com/rookie-xy/modules/proxy/forward/worker"
    "github.com/rookie-xy/hubble/job"
 //   "strings"
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
    sinceDB   = command.New( plugin.Flag, "client.sinceDB",    nil, "This option use to group" )
    pipe      = command.New( plugin.Flag, "pipeline.stream",  nil, "This option use to group" )
)

var commands = []command.Item{

    { batch,
      command.FILE,
      module.Proxy,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { pipe,
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

    { sinceDB,
      command.FILE,
      module.Proxy,
      command.SetObject,
      state.Enable,
      0,
      nil },
}

func New(log log.Log) module.Template {
    return &forward{
        log:    log,
        worker: worker.New(log),
        jobs:   job.New(log),
    }
}

func (r *forward) Init() {
    /*
    key := pipe.GetFlag() + "." + pipe.GetKey()
    fmt.Println("uuuuuuuuuuuuuuuuuuuuuuuuuuu: ", key)
    pipeline, err := factory.Pipeline(key, r.log, pipe.GetValue())
    if err != nil {
        fmt.Println("pipeline error ", err)
        return
    } else {
        r.queue = pipeline
    }

    name := client.GetKey()
    name = name[strings.LastIndex(name, ".") + 1:]

    register.Queue(name, pipeline)

    key = client.GetFlag() + "." + client.GetKey()
    client, err := factory.Client(key, r.log, client.GetValue())
    if err != nil {
        fmt.Println("client error ", err)
        return
    }

    key = sinceDB.GetFlag() + "." + sinceDB.GetKey()
    sinceDB, err := factory.Client(key, r.log, sinceDB.GetValue())
    if err != nil {
        fmt.Println("sinceDB error ", err)
        return
    }

    r.worker.Init(client, sinceDB)
    */

    return
}

func (f *forward) Main() {
    fmt.Println("Start proxy forward module ...")
/*
    for {
        event, err := f.queue.Dequeue(10)
        switch err {

        }

        f.worker.Q = event.(pipeline.Queue)
        f.jobs.Start(f.worker)
    }
*/
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
