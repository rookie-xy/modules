package reverse
/*
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

type reverse struct {
    log       log.Log
    server    proxy.Reverse
    pipeline  queue.Queue
}

var (
    pipeline  = command.New( plugin.Flag, "pipeline.stream",  nil, "This option use to group" )
    server    = command.New( plugin.Flag, "server.sincedb",    nil, "This option use to group" )
)

var commands = []command.Item{

    { pipeline,
      command.FILE,
      module.Proxy,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { server,
      command.FILE,
      module.Proxy,
      command.SetObject,
      state.Enable,
      0,
      nil },

}

func New(log log.Log) module.Template {
    return &reverse{
        log: log,
    }
}

func (r *reverse) Init() {
    key := pipeline.GetFlag() + "." + pipeline.GetKey()
    pipeline, err := factory.Pipeline(key, r.log, pipeline.GetValue())
    if err != nil {
        fmt.Println("pipeline error", err)
        return
    } else {
        r.pipeline = pipeline
    }

    register.Clones(server.GetKey(), pipeline)

    key = server.GetFlag() + "." + server.GetKey()
    if server, err := factory.Server(key, r.log, server.GetValue()); err != nil {
        fmt.Println("server error", err)
        return
    } else {
        r.server = server
    }

    register.Service(server.GetKey(), r.server)

    return
}

func (r *reverse) Main() {
    if r.server == nil || r.pipeline == nil {
        return
    }

    fmt.Println("Start proxy server module ...")

    for {
        event, status := r.pipeline.Dequeue(10)

        switch status {

        case state.Ignore:
            continue

        case state.Busy:
            //TODO sleep
        }

        r.server.Post(event)
    }
}

func (r *reverse) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Proxy, Name, commands, New)
}
*/
