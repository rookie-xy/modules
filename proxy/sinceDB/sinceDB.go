package sinceDB

import (
//    "fmt"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/pipeline"
    "github.com/rookie-xy/hubble/plugin"
    "github.com/rookie-xy/hubble/adapter"
    "fmt"
    "github.com/rookie-xy/modules/proxy/sinceDB/utils"
)

const Name  = "sinceDB"

type sincedb struct {
    log       log.Log

    pipeline  pipeline.Queue
    client    adapter.SinceDB

    batch     int
}

var (
    Pipeline  = command.New( plugin.Flag, "pipeline.channel",  nil, "This option use to group" )
    batch     = command.New( module.Flag, "batch",    64, "This option use to group" )
    client    = command.New( plugin.Flag, "client.sinceDB",    nil, "This option use to group" )
)

var commands = []command.Item{

    { Pipeline,
      command.FILE,
      module.Proxy,
      Name,
      command.SetObject,
      nil },

    { batch,
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

}

func New(l log.Log) module.Template {
    return &sincedb{
        log: l,
    }
}

func (s *sincedb) Init() {
    key := Pipeline.GetFlag() + "." + Pipeline.GetKey()
    pipeline, err := factory.Pipeline(key, s.log, Pipeline.GetValue())
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

    if value := batch.GetValue(); value != nil {
        if batch, err := value.GetInt(); err != nil {
            s.batch = batch
        }
    }

    return
}

func (s *sincedb) Main() {
    if s.client == nil || s.pipeline == nil {
        return
    }

    fmt.Println("Start proxy sinceDB module ...")

    for {
        events, err := s.pipeline.Dequeues(s.batch)
        switch err {

        case pipeline.ErrClosed:
            fmt.Println("sinceDB proxy close ...")
            return

        case pipeline.ErrEmpty:
        default:
        }

        if events != nil {
            if err := s.client.Senders(events); err != nil {
                if err := utils.Recall(events, s.pipeline); err != nil {
                    fmt.Println("recall error ", err)
                    return
                }
            }
        }
    }
}

func (s *sincedb) Exit(code int) {
    s.pipeline.Close()
	s.client.Close()
}

func init() {
    register.Module(module.Proxy, Name, commands, New)
}

        /*
        if events != nil {
            for _, event := range events {
            	// why? event is nil?
                if event != nil {
                    fileEvent := adapter.ToFileEvent(event)
                    fmt.Println("DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD ", fileEvent.GetFooter().Offset)
                }
            }
        }
        */
