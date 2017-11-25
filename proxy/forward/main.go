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
    log        log.Log

    pipeline   pipeline.Queue
    jobs      *job.Jobs
    done       chan struct{}

    client    *command.Command
    sinceDB   *command.Command

    name       string
}

var (
    client    = command.New( plugin.Flag, "client.kafka",    nil, "This option use to group" )
    sinceDB   = command.New( plugin.Flag, "client.sinceDB",    nil, "This option use to group" )
    Pipeline  = command.New( plugin.Flag, "pipeline.stream",  nil, "This option use to group" )
)

var commands = []command.Item{

    { Pipeline,
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
        log:  log,
        jobs: job.New(log),
        done: make(chan struct{}),
    }
}

func (f *forward) Init() {
    fmt.Println("Initialization forward component for proxy")
    key := Pipeline.GetFlag() + "." + Pipeline.GetKey()
    pipeline, err := factory.Pipeline(key, f.log, Pipeline.GetValue())
    if err != nil {
        fmt.Println("pipeline error ", err)
        return
    } else {
        f.pipeline = pipeline
    }

    name := client.GetKey()
    name = name[strings.LastIndex(name, ".") + 1:]

    register.Queue(name, pipeline)

    f.name = name
    f.client = command.New(
        client.GetFlag(),
        client.GetKey(),
       nil,
       "")

    f.sinceDB = command.New(
        sinceDB.GetFlag(),
        sinceDB.GetKey(),
       nil,
       "")
}

func (f *forward) Main() {
    fmt.Println("Start proxy forward module ... ", f.name, f.client.GetKey())
    defer func(jobs *job.Jobs) {
        jobs.WaitForCompletion()
        close(f.done)
    }(f.jobs)

    for {
        event, err := f.pipeline.Dequeue()
        switch err {
        case pipeline.ErrClosed:
        	fmt.Println("forwarder close ...")
        	return

        case pipeline.ErrEmpty:
        default:
            fmt.Println("Unknown error")
        }

        worker := worker.New(f.log)
        if err := worker.Init(f.client, f.sinceDB, event); err != nil {
            fmt.Println(err)
            return
        }

        f.jobs.Start(worker)
    }
}

func (f *forward) Exit(code int) {
    defer func() {
        <-f.done
        fmt.Println("Forward proxy component have exit")
    }()

    f.pipeline.Close()
/*
    if length := f.jobs.Len(); length > 0 {
        f.jobs.Stop()
    }
*/
}

func init() {
    register.Module(module.Proxy, Name, commands, New)
}
