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

var (
    resource = &command.Meta{ "-r", "resource", "./usr/local/conf/worker.yaml", "Resource type" }
)

var commands = []command.Item{

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

func New() *file {
    return &file{}
}

func (r *file) Init() {
fmt.Println("fileffffffffffffffff inittttttttttttt", resource.Value.(string))
    // 初始化文件监视器，监控配置文件
    // 初始化文件解析器解析文件


    return
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
    instance.Register(Name, commands, New())
}
