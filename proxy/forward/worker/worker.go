package worker

import (
    "github.com/satori/go.uuid"

    "github.com/rookie-xy/hubble/log"
	"github.com/rookie-xy/hubble/proxy"
	"github.com/rookie-xy/hubble/output"
	"fmt"
	"github.com/rookie-xy/hubble/pipeline"
	"github.com/rookie-xy/hubble/adapter"
	"github.com/rookie-xy/hubble/command"
	"github.com/rookie-xy/hubble/event"
)

type Worker struct {
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
/*
func (w *Worker) Init(f proxy.Forward, o output.Output) error {
	w.client = f
	w.sinceDB = o
	return nil
}
*/

func (w *Worker) Init(client, sinceDB *command.Command, event event.Event) error {
	w.client = f
	w.sinceDB = o
	return nil
}

func (w *Worker) ID() uuid.UUID {
    return w.id
}

func (w *Worker) Run() error {
	handle := func(Q pipeline.Queue, client proxy.Forward, sinceDB output.Output) error {
		for {
			event, err := Q.Dequeue(10)
			switch err {

			default:
			}

            body := adapter.ToFileEvent(event).GetBody()

			fmt.Println("workerrrrrrrrrrrrrrrrrrrrrrrrrrrr ", string(body.GetContent()))

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
	}

	return handle(w.Q, w.client, w.sinceDB)
}

func (w *Worker) Stop() {
	w.sinceDB.Close()
	w.client.Close()
	return
}
