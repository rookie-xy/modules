package file

import (
    "fmt"
    "unsafe"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/log"
)

const Name  = "file"

type fileInput struct{
    log.Log
}

func New(log log.Log) module.Template {
    return &fileInput{
        Log: log,
    }
}

var (
    group   = &command.Meta{ "", "group",   "nginx", "This option use to group" }
    types   = &command.Meta{ "", "type",    "log",   "file type, this is use to find some question" }
    paths   = &command.Meta{ "", "paths",   nil,     "File path, its is manny option" }
    publish = &command.Meta{ "", "publish", nil,     "publish topic" }
    codec   = &command.Meta{ "", "codec",   nil,     "codec method" }
)

var commands = []command.Item{

    { group,
      command.FILE,
      module.Inputs,
      command.SetObject,
      unsafe.Offsetof(group.Value),
      nil },

    { types,
      command.FILE,
      module.Inputs,
      command.SetObject,
      unsafe.Offsetof(types.Value),
      nil },

    { paths,
      command.FILE,
      module.Inputs,
      command.SetObject,
      unsafe.Offsetof(paths.Value),
      nil },

    { publish,
      command.FILE,
      module.Inputs,
      command.SetObject,
      unsafe.Offsetof(publish.Value),
      nil },

    { codec,
      command.FILE,
      module.Inputs,
      command.SetObject,
      unsafe.Offsetof(codec.Value),
      nil },
}

func (r *fileInput) Init() {
    //利用group codec等,进行初始化

    fmt.Println("qqqqqqqqqqqqqqqqqqqqqqqqqqqqq", group.Value, types.Value, paths.Value, publish.Value, codec.Value)


    return
}

func (r *fileInput) Main() {
    // 编写主要业务逻辑
}

func (r *fileInput) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Inputs, Name, commands, New)
}
