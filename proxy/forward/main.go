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
)

const Name  = "forward"

type forward struct {
    log.Log
    level      Level

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

    key := Pipeline.GetFlag() + "." + Pipeline.GetKey()
    pipeline, err := factory.Pipeline(key, f.Log, Pipeline.GetValue())
    if err != nil {
        f.log(ERROR, Name +"; pipeline error ", err)
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
        if err := worker.Init(f.client, f.sinceDB, event, f.log); err != nil {
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
    register.Module(module.Proxys, Name, commands, New)
}
