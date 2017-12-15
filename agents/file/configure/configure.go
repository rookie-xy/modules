package configure

import (
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/command"
    "time"
    "github.com/rookie-xy/hubble/output"
    "github.com/rookie-xy/hubble/proxy"
)

type Configure struct {
    Group     string
    Type      string
    Paths     types.Value
    Excludes  types.Value
    Limit     uint64
    Expire    time.Duration

    Client    proxy.Forward
    Output    output.Output
    SinceDB   output.Output
}
