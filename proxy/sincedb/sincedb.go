package sincedb

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
    //"github.com/rookie-xy/hubble/event"
    //"github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/event"
    "github.com/rookie-xy/modules/proxy/forward/events"
)

const Name  = "sincedb"

type sincedb struct {
    log       log.Log
    pipeline  queue.Queue
    client    adapter.SinceDB
}

var (
    pipeline  = command.New( plugin.Flag, "pipeline.stream",  nil, "This option use to group" )
    client    = command.New( plugin.Flag, "client.sincedb",    nil, "This option use to group" )
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

func New(l log.Log) module.Template {
    return &sincedb{
        log: l,
    }
}

func (s *sincedb) Init() {
    key := pipeline.GetFlag() + "." + pipeline.GetKey()
    pipeline, err := factory.Pipeline(key, s.log, pipeline.GetValue())
    if err != nil {
        fmt.Println("pipeline error ", err)
        return
    } else {
        s.pipeline = pipeline
    }

    register.Queue(client.GetKey(), pipeline)

    key = client.GetFlag() + "." + client.GetKey()
    if client, err := factory.Client(key, s.log, client.GetValue()); err != nil {
        fmt.Println("client error ", err)
        return
    } else {
        s.client = adapter.FileSinceDB(client)
        register.Forword(key, client)
    }

    return
}

func (s *sincedb) Main() {
    if s.client == nil || s.pipeline == nil {
        return
    }

    fmt.Println("Start proxy forward module ...")

    for {
        event, status := s.pipeline.Dequeue(10)

        switch status {
/*
        case state.Ignore:
            continue
        case state.Busy:
            //TODO sleep
*/
        default:
        }

        if !s.client.Commit(event) {
            if events, err := s.client.Sender(); err != nil {
                if err := recall(events, s.pipeline); err != nil {
                    fmt.Println("recall error ", err)
                    return
                }
            }
        }
    }
}

func recall(events []event.Event, Q queue.Queue) error {
    for _, event := range events {
        if err := Q.Requeue(event); err != nil {
            return err
        }
    }

    return nil
}

func (s *sincedb) Exit(code int) {
	s.client.Close()
    // 退出
}

func init() {
    register.Module(module.Proxy, Name, commands, New)
}
