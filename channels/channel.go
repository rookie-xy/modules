package channels

import (
    "github.com/rookie-xy/worker/src/event"
    "github.com/rookie-xy/worker/src/state"
    "github.com/rookie-xy/worker/src/log"
)

type channel struct {
    log.Log
    pipline chan event.Event
}

func Create(log log.Log, size int) *channel {
    return &channel{
        Log: log,
        pipline: make(chan event.Event, size),
    }
}

// TODO 确定如何保证并发
func (r *channel) Clone() *channel {
    return r
}

func (r *channel) Push(e event.Event) int {
    r.pipline <- e
    return state.Ok
}

func (r *channel) Pull(size int) (event.Event, int) {
    event, open := <- r.pipline
    if !open {
        return nil, state.Done
    }

    return event, state.Ok
}
