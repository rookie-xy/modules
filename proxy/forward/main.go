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
)

const Name  = "forward"

type forward struct {
    log       log.Log
    client    proxy.Forward
    sincedb   output.Output
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

    var output output.Output

    fmt.Println("Start proxy forward module ...")

    for {
    	conn := output.Accept()

        handler := func(Q pipeline.Queue, forward proxy.Forward, sincedb output.Output) error {
            for {
                event, err := Q.Dequeue(10)

                switch err {

                default:
                }

                if err := forward.Sender(event); err != nil {
                    if err = Q.Requeue(event); err != nil {
                        fmt.Println("recall error ", err)
                        return err
                    }
                    continue
                }

                if err := sincedb.Sender(event); err != nil {
                    fmt.Println("sincedb sender error ", err)
                    return err
                }
            }
        }

        go handler(conn, f.client, f.sincedb)
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
