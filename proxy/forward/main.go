package forward

import (
    "strings"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
  . "github.com/rookie-xy/hubble/log/level"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/plugin"
    "github.com/rookie-xy/hubble/pipeline"
    "github.com/rookie-xy/hubble/job"

    "github.com/rookie-xy/modules/proxy/forward/worker"
    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/types/value"
    "github.com/rookie-xy/hubble/proxy"
    "github.com/rookie-xy/hubble/output"
)

const Name  = "forward"

type forward struct {
    log.Log
    level      Level

    pipeline   pipeline.Queue
    jobs      *job.Jobs
    done       chan struct{}

    client     proxy.Forward
    sinceDB    output.Output

    name       string
}

var (
    client    = command.New( plugin.Flag, "client.kafka",    nil, "This option use to group" )
    sinceDB   = command.New( plugin.Flag, "client.sinceDB",  nil, "This option use to group" )
    Pipeline  = command.New( plugin.Flag, "pipeline.stream", nil, "This option use to group" )
)

var commands = []command.Item{

    { Pipeline,
      command.FILE,
      module.Proxys,
      Name,
      command.SetObject,
      nil },

    { client,
      command.FILE,
      module.Proxys,
      Name,
      command.SetObject,
      nil },

    { sinceDB,
      command.FILE,
      module.Proxys,
      Name,
      command.SetObject,
      nil },
}

func New(log log.Log) module.Template {
    return &forward{
        Log:  log,
        level: adapter.ToLevelLog(log).Get(),
        jobs: job.New(log),
        done: make(chan struct{}),
    }
}

func (f *forward) Init() {
    f.log(DEBUG, Name +"; init component for proxy")

    if key, ok := plugin.Name(Pipeline.GetKey()); ok {
        pipeline, err := factory.Pipeline(key, f.Log, Pipeline.GetValue())
        if err != nil {
            f.log(ERROR, Name + "; pipeline error ", err)
            return
        } else {
            f.pipeline = pipeline
        }
    }

    name := client.GetKey()
    name = name[strings.LastIndex(name, ".") + 1:]

    register.Queue(name, f.pipeline)
    f.name = name

    if key, ok := plugin.Name(client.GetKey()); ok {
        var err error
        f.client, err = factory.Client(key, f.Log, client.GetValue())
        if err != nil {
        	f.log(ERROR, Name + "; client error ", err)
            return
        }
    }

    if key, ok := plugin.Domain(output.Name, "sinceDB"); ok {
        var err error
        f.sinceDB, err = factory.Output(key, f.Log, value.New(sinceDB.GetKey()))
        if err != nil {
         	f.log(ERROR, Name + "; sinceDB error: %s", err)
            return
        }
    }
}

func (f *forward) Main() {
    f.log(INFO, Name +"; run component for %s", f.name)

    defer func(jobs *job.Jobs) {
        jobs.WaitForCompletion()
        close(f.done)
    }(f.jobs)

    for {
        event, err := f.pipeline.Dequeue()
        switch err {
        case pipeline.ErrClosed:
        	f.log(INFO, Name +"; close for %s, %s", f.name, pipeline.ErrClosed)
        	return

        case pipeline.ErrEmpty:
            f.log(INFO, Name +"; empty for %s, %s", f.name, pipeline.ErrEmpty)
        default:
            f.log(WARN, Name +"; unknown queue event")
        }

        worker := worker.New(f.Log)
        if err := worker.Init(f.client, f.sinceDB, event); err != nil {
            f.log(WARN, Name + "; %s", err.Error())
            return
        }

        f.jobs.Start(worker)
    }
}

func (f *forward) Exit(code int) {
    defer func() {
        <-f.done
        f.log(DEBUG,"%s component have exit", Name)
    }()

    f.log(INFO,"Exit component for %s", Name)
    f.pipeline.Close()
/*
    if length := f.jobs.Len(); length > 0 {
        f.jobs.Stop()
    }
*/
}

func (f *forward) log(l Level, fmt string, args ...interface{}) {
    log.Print(f.Log, f.level, l, fmt, args...)
}

func init() {
    register.Component(module.Proxys, Name, commands, New)
}
