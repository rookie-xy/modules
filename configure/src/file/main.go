package file

import (
    "unsafe"
    "github.com/rookie-xy/worker/src/cycle"
    "github.com/rookie-xy/worker/src/instance"
)

const Name  = "file"

var (
    resource = command.Meta{ "-r", "resource", "./usr/local/conf/worker.yaml", "Resource type" }
    format   = command.Meta{ "-f", "format", "yaml", "Configure file format" }
)

var commands = []command.Item{

    { resource,
      command.LINE,
      module.GLOBEL,
      configure.SetString,
      unsafe.Offsetof(resource.Value),
      nil },

    { format,
      command.LINE,
      module.GLOBEL,
      configure.SetString,
      unsafe.Offsetof(format.Value),
      nil },
}

type File struct {
    log.Log
}

func New() *File {
    return &File{}
}

func (r *File) Init() {

    // 初始化文件监视器，监控配置文件
    // 初始化文件解析器解析文件


    return
}

func (r *File) Main() {
    // 发现文件变更，通知给其他模块

    return
}

func (r *File) Exit() {
    for _, module := range r.children {
        module.Exit()
    }

    //r.cycle.Quit()
    return
}

func init() {
    instance.Register(Name, commands, New())
}
