package file

import (
    "unsafe"
    //"hubble/modules/inputs/src/file/src"
    "hubble/src/command"
    "hubble/src/configure"
    "hubble/src/module"
"hubble/src/register"
)

const Name  = "file"

type fileInput struct{
}

func New() *fileInput {
    return &fileInput{}
}

var (
    group   = command.Name{ "", "group",   "nginx",     "This option use to group" }
    types   = command.Name{ "", "type",    "log",       "file type, this is use to find some question" }
    paths   = command.Name{ "", "paths",   array.New(), "File path, its is manny option" }
    publish = command.Name{ "", "publish", array.New(), "publish topic" }
    codec   = command.Name{ "", "codec",   codec.New(), "codec method" }
)

var fileInputCommands = []command.Command{

    { group,
      command.FILE,
      module.LOCAL,
      configure.SetString,
      unsafe.Offsetof(group.Value),
      nil },

    { types,
      command.FILE,
      module.LOCAL,
      configure.SetString,
      unsafe.Offsetof(types.Value),
      nil },

    { paths,
      command.FILE,
      module.LOCAL,
      configure.SetArray,
      unsafe.Offsetof(paths.Value),
      nil },

    { publish,
      command.FILE,
      module.LOCAL,
      configure.SetArray,
      unsafe.Offsetof(publish.Value),
      nil },

    { codec,
      command.FILE,
      module.LOCAL,
      configure.SetCodec,
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
    r.cycle.Stop()
    // 退出
}

func init() {
    register.GetInstance().Module(Name, New(), fileInputCommands)
}

