package configure

import (
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/hubble/input"
    "time"
)

type Configure struct {
    Group     string
    Type      string
    Paths     types.Value
    Excludes  types.Value
    Limit     uint64
    Expire    time.Duration

    Input     input.Input
    Codec     codec.Codec

    Client    *command.Command
    Output    *command.Command
    SinceDB   *command.Command
}
