package configure

import (
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/command"
    "time"
)

type Configure struct {
    Group     string
    Type      string
    Paths     types.Value
    Excludes  types.Value
    Limit     uint64
    Expire    time.Duration

    Client    *command.Command
    Output    *command.Command
    SinceDB   *command.Command
}
