package remote

import (
    "fmt"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/log"
)

const Name = "remote"

var (
    client    = command.New("-c", "client",    "etcd", "If you want to ")
    parameter = command.New("-p", "parameter", "10.0.1.2", "If you want to ")
)

var commands = []command.Item{

    { client,
      command.LINE,
      module.Configure,
      Name,
      command.SetObject,
      nil },


    { parameter,
      command.LINE,
      module.Configure,
      Name,
      command.SetObject,
      nil },

}

type remote struct {
    log.Log
}

func New(log log.Log) module.Template {
    return &remote{
        Log: log,
    }
}

func (r *remote) Init() {
    if v := client.GetValue(); v != nil {
        fmt.Println(v.GetString())
    }

    if v1 := parameter.GetValue(); v1 != nil {
        fmt.Println(v1.GetString())
    }
    // 初始化文件解析器解析文件

    return
}

func (r *remote) Main() {
    return
}

func (r *remote) Exit(code int) {
    //r.cycle.Quit()
    return
}

func init() {
    register.Module(module.Configure, Name, commands, New)
}
