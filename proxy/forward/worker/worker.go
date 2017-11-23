package worker

import (
	"sync"

    "github.com/satori/go.uuid"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/log"
	"github.com/rookie-xy/hubble/proxy"
	"github.com/rookie-xy/hubble/output"
	"fmt"
	"github.com/rookie-xy/hubble/pipeline"
	"github.com/rookie-xy/hubble/adapter"
	"github.com/rookie-xy/hubble/command"
	"github.com/rookie-xy/hubble/event"
	"github.com/rookie-xy/hubble/types/value"
)

type Worker struct {
    sync.Mutex
    id        uuid.UUID

    Q         pipeline.Queue
    client    proxy.Forward
    sinceDB   output.Output

    log  log.Log
}

func New(log log.Log) *Worker {
	return &Worker{
		log: log,
		id: uuid.NewV4(),
	}
}

func (w *Worker) Init(pclient, psinceDB *command.Command, event event.Event) error {
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

	w.client = client
    w.sinceDB = sinceDB
	w.Q = adapter.ToPipelineEvent(event)

	return nil
}

func (w *Worker) ID() uuid.UUID {
    return w.id
}

func (w *Worker) Run() error {
	fmt.Println("WOZAIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIII")

	defer func() {
	    w.sinceDB.Close()
        w.client.Close()
	}()

	handle := func(Q pipeline.Queue, client proxy.Forward, sinceDB output.Output) error {

		keep := true

		for {
			event, err := Q.Dequeue()
			switch err {

			case pipeline.ErrClosed:
				fmt.Println("CHANNELLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL CLOSE")
				keep = false

			case pipeline.ErrEmpty:

			default:
			}

			if !keep {
				break
			}

			if err := client.Sender(event); err != nil {
				if err = Q.Requeue(event); err != nil {
                    //w.log.Print("aaa")
					fmt.Println("recall error ", err)
					return err
				}
				continue
			}

			if err := sinceDB.Sender(event); err != nil {
				fmt.Println("sinceDB sender error ", err)
				return err
			}
		}

	    return nil
	}

	return handle(w.Q, w.client, w.sinceDB)
}

func (w *Worker) Stop() {
    w.Q.Close()
}
