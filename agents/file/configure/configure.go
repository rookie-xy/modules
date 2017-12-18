package configure

import (
    "github.com/rookie-xy/hubble/types"
    "time"
    "github.com/rookie-xy/hubble/output"
)

type Configure struct {
    Group     string
    Type      string
    Paths     types.Value
    Excludes  types.Value
    Limit     uint64
    Expire    time.Duration

    Client    bool
    Output    output.Output
    SinceDB   output.Output
}
