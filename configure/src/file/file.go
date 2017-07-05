package file

import (
    "unsafe"
    "github.com/rookie-xy/worker/src/instance"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
        "fmt"
)

const Name  = "file"

var Event chan string = make(chan string)

var (
    //format   = &command.Meta{ "-f", "format", "yaml", "Configure file format" }
    format   = &command.Meta{ "-f", "file", "yaml", "Configure file format" }
    resource = &command.Meta{ "-p", "path", "./usr/local/conf/worker.yaml", "Resource type" }
)

var commands = []command.Item{

    { format,
      command.LINE,
      module.GLOBEL,
      command.SetObject,
      unsafe.Offsetof(format.Value),
      nil },

    { resource,
      command.LINE,
      module.GLOBEL,
      command.SetObject,
      unsafe.Offsetof(resource.Value),
      nil },

}

type file struct {
    log.Log
}

func New(log log.Log) *file {
    return &file{
            Log: log,
    }
}

func Init(log log.Log) module.Template {
    file := New(log)

    //fmt.Println("fileffffffffffffffff inittttttttttttt", resource.Value.(string))
    // 判断文件是否存在，可读性
    // 初始化文件监视器，监控配置文件
    // 初始化解析器

    return file
}

func (r *file) Main() {
fmt.Println("fileffffffffffffffff mainnnnnnnnnnnnnnnnn")
    // 发现文件变更，通知给其他模块

    return
}

func (r *file) Exit() {
    //r.cycle.Quit()
    return
}

func init() {
    instance.Register(Name, Name, commands, Init)
}
