package configure

import (
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/command"
)

type Configure struct {
    Group     string
    Type      string
    Paths     types.Value
    Excludes  types.Value
    Limit     uint64

    Source   *command.Command
    Codec    *command.Command
    Client   *command.Command
}
