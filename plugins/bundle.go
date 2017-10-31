package plugins

import (
    "fmt"

    "github.com/rookie-xy/hubble/plugin"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/register"

  _ "github.com/rookie-xy/plugins/codec"
  _ "github.com/rookie-xy/plugins/client"
  _ "github.com/rookie-xy/plugins/pipeline"
  _ "github.com/rookie-xy/plugins/valve"
  _ "github.com/rookie-xy/plugins/input"
  _ "github.com/rookie-xy/plugins/output"
)

const Name = "dso"

var (
    dso = command.New("-so", "plugin", paths, "You can dynamically load the plugin in DSO mode")
)

var commands = []command.Item{

    { dso,
      command.LINE,
      module.Plugins,
      Name,
      command.SetObject,
      nil },

}

var paths = [...]string{}

type plugins struct {
}

func (r *plugins) Init() {
    if len(paths) > 0 {
        fmt.Println("EXPERIMENTAL: loadable plugin support is experimental")
    }

    for _, path := range paths {
        fmt.Println("loading plugin bundle: %v", path)

        if err := plugin.Load(path); err != nil {
            fmt.Println(err)
            return
        }

        return
    }
}

func (r *plugins) Main() {
    // nothing
}

func (r *plugins) Exit(code int) {
    // nothing
}

func init() {
    register.Module(module.Worker, module.Plugins, commands, nil)
}
