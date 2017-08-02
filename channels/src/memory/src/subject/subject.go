package subject

import (
    "fmt"
    "github.com/rookie-xy/worker/src/observer"
    "github.com/rookie-xy/worker/src/log"
)

type Subject struct {
    log.Log
    observers  []observer.Observer
    data       string
}

func New(log log.Log) *Subject {
    return &Subject{
        Log: log,
    }
}

func (r *Subject) Attach(o observer.Observer) {
    if o != nil {
        r.observers = append(r.observers, o)
        return
    }

    fmt.Println("attach error")
    return
}

func (r *Subject) Notify() {
    if r.data == "" {
        return
    }

    fmt.Println(r.data)
}
