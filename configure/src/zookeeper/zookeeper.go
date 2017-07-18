package zookeeper

import (
    "unsafe"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/log"
    "fmt"
)

const Name  = "zookeeper"

var (
    address = &command.Meta{ "-a", "address", "192.168.1.1:2181", "Resource type" }
)

var commands = []command.Item{

    { address,
      command.LINE,
      module.Configure,
      command.SetObject,
      unsafe.Offsetof(address.Value),
      nil },

}

type zookeeper struct {
    log.Log
}

func New(log log.Log) module.Template {
    return &zookeeper{
        Log: log,
    }
}

func (r *zookeeper) Init() {
fmt.Println("zookeeperffffffffffffffff inittttttttttttt", address.Value.(string))
    // 初始化zkClient
    // 初始化文件解析器解析文件

    return
}

func (r *zookeeper) Main() {
fmt.Println("zookeeperffffffffffffffff mainnnnnnnnnnnnnnnnn")
    // 从zk拉取配置
    // 解析配置
    // 吐出数据

    return
}

func (r *zookeeper) Exit(code int) {
    //r.cycle.Quit()
    return
}

func init() {
    register.Module(module.Configure, Name, commands, New)
}