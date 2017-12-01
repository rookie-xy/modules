package worker

import (
	"sync"

    "github.com/satori/go.uuid"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/log"
	"github.com/rookie-xy/hubble/proxy"
	"github.com/rookie-xy/hubble/output"
	"github.com/rookie-xy/hubble/pipeline"
	"github.com/rookie-xy/hubble/adapter"
	"github.com/rookie-xy/hubble/command"
	"github.com/rookie-xy/hubble/event"
	"github.com/rookie-xy/hubble/types/value"
  . "github.com/rookie-xy/hubble/log/level"
)

type Worker struct {
    sync.Mutex
    id        uuid.UUID
    name      string

    Q         pipeline.Queue
    client    proxy.Forward
    sinceDB   output.Output

    log  log.Log
    logf  log.Factory
}

func New(log log.Log) *Worker {
	return &Worker{
		log: log,
		id: uuid.NewV4(),
	}
}

func (w *Worker) Init(pclient, psinceDB *command.Command, event event.Event, logf log.Factory) error {
	w.name = pclient.GetKey()
    key := pclient.GetFlag() + "." + pclient.GetKey()
    client, err := factory.Client(key, w.log, pclient.GetValue())
    if err != nil {
        return err
    }

    key = psinceDB.GetFlag() + "." + output.Name + "." + "sinceDB"
    sinceDB, err := factory.Output(key, w.log, value.New(psinceDB.GetKey()))
    if err != nil {
        return err
    } else {
        sinceDB.Sender(nil)
    }

    w.logf = logf
	w.client = client
    w.sinceDB = sinceDB
	w.Q = adapter.ToPipelineEvent(event)

	return nil
}

func (w *Worker) ID() uuid.UUID {
    return w.id
}

func (w *Worker) Run() error {
    w.logf(DEBUG,"worker running for %s", w.name)

	defer func() {
	    w.sinceDB.Close()
        w.client.Close()
	}()

	handle := func(Q pipeline.Queue, client proxy.Forward, sinceDB output.Output, log log.Factory) error {

		keep := true
		for {
			event, err := Q.Dequeue()
			switch err {

			case pipeline.ErrClosed:
				keep = false
				log(INFO, "forward worker; close %s", pipeline.ErrClosed)

			case pipeline.ErrEmpty:
				log(INFO, "forward worker; empty %s", pipeline.ErrEmpty)
			default:
				log(WARN, "forward worker; unknown queue event")
			}

			if !keep {
				break
			}

			if err := client.Sender(event); err != nil {
				if err = Q.Requeue(event); err != nil {
					log(ERROR,"worker; recall error: %s", err)
					return err
				}
				continue
			}

			if err := sinceDB.Sender(event); err != nil {
                log(ERROR, "worker; sinceDB sender error: %s", err)
				return err
			}
		}

	    return nil
	}

	return handle(w.Q, w.client, w.sinceDB, w.logf)
}

func (w *Worker) Stop() {
	w.logf(DEBUG, "worker stop for %v", w.id)
    w.Q.Close()
}
