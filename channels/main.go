package channels


import (
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
)

const Name = module.Channels

var (
    channels = &command.Meta{ "", Name, nil, "inputs may be many" }
)

var commands = []command.Item{

    { channels,
      command.FILE,
      Name,
      nil,
      0,
      nil },

}
