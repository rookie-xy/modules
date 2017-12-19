package worker

import (
	"sync"

    "github.com/satori/go.uuid"
    "github.com/rookie-xy/hubble/log"
	"github.com/rookie-xy/hubble/proxy"
	"github.com/rookie-xy/hubble/output"
	"github.com/rookie-xy/hubble/pipeline"
	"github.com/rookie-xy/hubble/adapter"
	"github.com/rookie-xy/hubble/event"
  . "github.com/rookie-xy/hubble/log/level"
	"github.com/rookie-xy/hubble/prototype"
)

type Worker struct {
    sync.Mutex
    id        uuid.UUID
    name      string

    Q         pipeline.Queue
    client    proxy.Forward
    sinceDB   proxy.Forward

    level     Level
    log.Log
}

func New(log log.Log) *Worker {
	return &Worker{
		Log: log,
		level: adapter.ToLevelLog(log).Get(),
		id: uuid.NewV4(),
	}
}

func (w *Worker) Init(client proxy.Forward, sinceDB output.Output, event event.Event) error {
    if sinceDB, err := sinceDB.New(); err != nil {
    	return err
	} else {
		w.sinceDB = sinceDB
	}

    w.client  = prototype.Forward(client)
	w.Q       = adapter.ToPipelineEvent(event)
	return nil
}

func (w *Worker) ID() uuid.UUID {
    return w.id
}

func (w *Worker) Run() error {
    w.log(DEBUG,"worker running for %s", w.name)

	defer func() {
	    w.sinceDB.Close()
        w.client.Close()
	}()

	handle := func(Q pipeline.Queue, client, sinceDB proxy.Forward) error {

		keep := true
		for {
			event, err := Q.Dequeue()
			switch err {

			case pipeline.ErrClosed:
				keep = false
				w.log(INFO, "forward worker; close %s", pipeline.ErrClosed)

			case pipeline.ErrEmpty:
				w.log(INFO, "forward worker; empty %s", pipeline.ErrEmpty)
			default:
				w.log(WARN, "forward worker; unknown queue event")
			}

			if !keep {
				break
			}

			if err := client.Sender(event); err != nil {
				if err = Q.Requeue(event); err != nil {
				    w.log(ERROR,"worker; recall error: %s", err)
					return err
				}
				continue
			}

			if err := sinceDB.Sender(event); err != nil {
                w.log(ERROR, "worker; sinceDB sender error: %s", err)
				return err
			}
		}

	    return nil
	}

	return handle(w.Q, w.client, w.sinceDB)
}

func (w *Worker) Stop() {
	w.log(DEBUG, "worker stop for %v", w.id)
    w.Q.Close()
}

func (w *Worker) log(l Level, fmt string, args ...interface{}) {
    log.Print(w.Log, w.level, l, fmt, args...)
}

	/*
	w.name = pclient.GetKey()
    key := pclient.GetFlag() + "." + pclient.GetKey()
    client, err := factory.Client(key, w.Log, pclient.GetValue())
    if err != nil {
        return err
    }

    key = psinceDB.GetFlag() + "." + output.Name + "." + "sinceDB"
    sinceDB, err := factory.Output(key, w.Log, value.New(psinceDB.GetKey()))
    if err != nil {
        return err
    } else {
        sinceDB.Sender(nil)
    }
	*/
