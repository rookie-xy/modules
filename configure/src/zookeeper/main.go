package zookeeper

import (
    "unsafe"
    "github.com/rookie-xy/worker/src/instance"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "fmt"
)

const Name  = "zookeeper"

var (
    resource = &command.Meta{ "-r", "resource", "192.168.1.1:2181", "Resource type" }
)

var commands = []command.Item{

    { resource,
      command.LINE,
      module.GLOBEL,
      command.SetObject,
      unsafe.Offsetof(resource.Value),
      nil },

}

type zookeeper struct {
    log.Log
}

func New() *zookeeper {
    return &zookeeper{}
}

func (r *zookeeper) Init() {
fmt.Println("zookeeperffffffffffffffff inittttttttttttt", resource.Value.(string))
    // 初始化文件监视器，监控配置文件
    // 初始化文件解析器解析文件


    return
}

func (r *zookeeper) Main() {
fmt.Println("zookeeperffffffffffffffff mainnnnnnnnnnnnnnnnn")
    // 发现文件变更，通知给其他模块

    return
}

func (r *zookeeper) Exit() {
    //r.cycle.Quit()
    return
}

func init() {
    instance.Register(Name, commands, New())
}
