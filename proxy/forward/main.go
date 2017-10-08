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
    queue "github.com/rookie-xy/hubble/pipeline"
    "github.com/rookie-xy/hubble/plugin"
    "github.com/rookie-xy/modules/proxy/forward/events"
    "github.com/rookie-xy/hubble/event"
    "github.com/rookie-xy/hubble/adapter"
)

const Name  = "forward"

type forward struct {
    log       log.Log
    pipeline  queue.Queue
    client    proxy.Forward
    events   *events.Batch
}

var (
    pipeline  = command.New( plugin.Flag, "pipeline.stream",  nil, "This option use to group" )
    client    = command.New( plugin.Flag, "client.elasticsearch",    nil, "This option use to group" )
    batch     = command.New( module.Flag, "events",    nil, "This option use to group" )
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

    { batch,
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
    key := pipeline.GetFlag() + "." + pipeline.GetKey()
    pipeline, err := factory.Pipeline(key, r.log, pipeline.GetValue())
    if err != nil {
        fmt.Println("pipeline error ", err)
        return
    } else {
        r.pipeline = pipeline
    }

    register.Queue(client.GetKey(), pipeline)

    key = client.GetFlag() + "." + client.GetKey()
    if client, err := factory.Client(key, r.log, client.GetValue()); err != nil {
        fmt.Println("client error ", err)
        return
    } else {
        r.client = client
        register.Forword(key, client)
    }

    events := events.New(r.log)
    if err := events.Init(batch.GetValue()); err != nil {
        fmt.Println("events init error ", err)
    	return
    } else {
        r.events = events
    }

    return
}

func (r *forward) Main() {
    if r.client == nil || r.pipeline == nil {
        return
    }

    batch := false
    if r.events.Enable() {
    	batch = true
	}

    fmt.Println("Start proxy forward module ...")

    for {
        event, status := r.pipeline.Dequeue(10)

        switch status {

        case state.Ignore:
            continue
        case state.Busy:
            //TODO sleep
        }

        if batch {
            r.events.Put(event)
            event = r.events
        }

        if err := r.client.Sender(event, batch); err != nil {
            if err = r.recall(event, r.pipeline, batch); err != nil {
                fmt.Println("recall error ", err)
            }
        }
    }
}

func (r *forward) recall(e event.Event, Q queue.Queue, batch bool) error {
    if batch {
        events := adapter.ToEvents(e)
        for _, event := range events.Batch() {
            if err := Q.Requeue(event); err != nil {
                return err
            }
        }

        return nil
    }

    return Q.Requeue(e)
}

func (r *forward) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Proxy, Name, commands, New)
}
