package file

import (
    "unsafe"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/instance"
"github.com/rookie-xy/worker/src/module"
)

const Name  = "file"

type fileInput struct{
}

func New() *fileInput {
    return &fileInput{}
}

var (
    group   = command.Meta{ "", "group",   "nginx",     "This option use to group" }
    types   = command.Meta{ "", "type",    "log",       "file type, this is use to find some question" }
    paths   = command.Meta{ "", "paths",   nil, "File path, its is manny option" }
    publish = command.Meta{ "", "publish", nil, "publish topic" }
    codec   = command.Meta{ "", "codec",   nil, "codec method" }
)

var commands = []command.Item{

    { group,
      command.FILE,
      module.LOCAL,
      command.SetObject,
      unsafe.Offsetof(group.Value),
      nil },

    { types,
      command.FILE,
      module.LOCAL,
      command.SetObject,
      unsafe.Offsetof(types.Value),
      nil },

    { paths,
      command.FILE,
      module.LOCAL,
      command.SetObject,
      unsafe.Offsetof(paths.Value),
      nil },

    { publish,
      command.FILE,
      module.LOCAL,
      command.SetObject,
      unsafe.Offsetof(publish.Value),
      nil },

    { codec,
      command.FILE,
      module.LOCAL,
      command.SetObject,
      unsafe.Offsetof(codec.Value),
      nil },
}

func (r *fileInput) Init() {
    //利用group codec等,进行初始化
    if group.Value != nil {
    }
}

func (r *fileInput) Main() {
    // 编写主要业务逻辑
}

func (r *fileInput) Exit() {
    // 退出
}

func init() {
    instance.Register(Name, commands, New)
}

