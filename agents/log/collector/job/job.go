package job

import (
    "github.com/rookie-xy/hubble/scanner"
    "github.com/rookie-xy/hubble/log"
)

type Job struct {
    log log.Log
}

func (j *Job) New(log log.Log) *Job {
    return &Job{
        log: log,
    }
}

func (c *job) ID() uuid.UUID {
    return 1
}

func (c *job) Run() error {
    scanner := scanner.New(nil)
    scanner.Split(c.codec.Decode)

    for scanner.Scan() {
        value := scanner.Value()
        value.GetString()
    }

    return nil
}

func (c *job) Stop() {
    return
}
