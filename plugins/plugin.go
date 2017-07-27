package plugins

import (
    "fmt"
    "unsafe"

    "github.com/rookie-xy/worker/src/plugin"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/register"

  _ "github.com/rookie-xy/plugins/codec"
)

var (
    dso = &command.Meta{ "-so", "plugin", paths, "You can dynamically load the plugin in DSO mode" }
)

var commands = []command.Item{

    { dso,
      command.LINE,
      module.Plugins,
      command.SetObject,
      unsafe.Offsetof(dso.Value),
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
    register.Module(module.Worker, plugin.Name, commands, nil)
}
