package outputs

import (
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
)

const Name = module.Inputs

var (
    outputs = &command.Meta{ "", Name, nil, "inputs may be many" }
)

var commands = []command.Item{

    { outputs,
      command.FILE,
      Name,
      nil,
      0,
      nil },

}
