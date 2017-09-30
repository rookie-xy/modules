package finder

import "github.com/rookie-xy/hubble/log"

type Finder struct {
    log  log.Log
}

func New(log log.Log) *Finder {
    return &Finder{
        log: log,
    }
}