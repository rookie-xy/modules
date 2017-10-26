package configure

import (
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/codec"
//    "github.com/rookie-xy/hubble/proxy"
    "github.com/rookie-xy/hubble/input"
//    "github.com/rookie-xy/hubble/output"
    "github.com/rookie-xy/hubble/proxy"
)

type Configure struct {
    Group     string
    Type      string
    Paths     types.Value
    Excludes  types.Value
    Limit     uint64

    Client    bool

    Input     input.Input
    Codec     codec.Codec
    Output    proxy.Forward

    Sincedb   command.Command
}
