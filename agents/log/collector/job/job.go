package job

import (
    "github.com/rookie-xy/hubble/scanner"
    "github.com/rookie-xy/hubble/log"
        "github.com/rookie-xy/hubble/codec"
    "github.com/satori/go.uuid"
)

type Job struct {
    log log.Log
    codec codec.Codec
}

func (j *Job) New(log log.Log) *Job {
    return &Job{
        log: log,
    }
}

func (c *Job) ID() uuid.UUID {
    return uuid.UUID{}
}

func (j *Job) Run() error {
    scanner := scanner.New(nil)
    scanner.Split(j.codec.Decode)

    for scanner.Scan() {
        value := scanner.Value()
        value.GetString()
    }

    return nil
}

func (j *Job) Stop() {
    return
}
