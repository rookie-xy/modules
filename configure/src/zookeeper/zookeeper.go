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
    format   = &command.Meta{ "-f", "format", "json", "Configure file format" }
    resource = &command.Meta{ "-r", "resource", "192.168.1.1:2181", "Resource type" }
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

type zookeeper struct {
    log.Log
}

func New() *zookeeper {
    return &zookeeper{}
}

func (r *zookeeper) Init(name string) module.Template {
fmt.Println("zookeeperffffffffffffffff inittttttttttttt", resource.Value.(string))
    // 初始化zkClient
    // 初始化文件解析器解析文件

    zk := New()


    return zk
}

func (r *zookeeper) Main() {
fmt.Println("zookeeperffffffffffffffff mainnnnnnnnnnnnnnnnn")
    // 从zk拉取配置
    // 解析配置
    // 吐出数据

    return
}

func (r *zookeeper) Exit() {
    //r.cycle.Quit()
    return
}

func init() {
    instance.Register(Name, Name, commands, New())
}
