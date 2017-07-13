package file

import (
    "unsafe"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/configure"
    "github.com/rookie-xy/worker/src/log"
    "fmt"
)

const Name  = "file"

var (
    path = &command.Meta{ "-p", "path", "./usr/local/conf/worker.yaml", "Resource type" }
)

var commands = []command.Item{

    { path,
      command.LINE,
      module.GLOBEL,
      command.SetObject,
      unsafe.Offsetof(path.Value),
      nil },

}

type file struct {
    log.Log
}

func New(log log.Log) module.Module {
    return &file{
        Log: log,
    }
}

func (r *file) Init() {
    //fmt.Println("fileffffffffffffffff inittttttttttttt", resource.Value.(string))
    // 判断文件是否存在，可读性
    // 初始化文件监视器，监控配置文件
    // 初始化解析器

    return
}

func (r *file) Main() {
fmt.Println("fileffffffffffffffff mainnnnnnnnnnnnnnnnn")
    // 发现文件变更，通知给其他模块

    configure.Event

    return
}

func (r *file) Exit() {
    //r.cycle.Quit()
    return
}

func init() {
    instance.Register(configure.Name, Name, commands, New)
}
