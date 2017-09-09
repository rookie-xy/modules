package stdin

import (
    "fmt"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/state"
    "github.com/rookie-xy/hubble/plugin"

    "github.com/rookie-xy/modules/agents/log/collector"
)

const Name  = "stdin"

type stdin struct{
    log.Log
}

func New(log log.Log) module.Template {
    return &stdin{
        Log: log,
    }
}

var (
    group   = command.New( module.Flag, "group",   "nginx", "This option use to group" )
    types   = command.New( module.Flag, "type",    "log",   "file type, this is use to find some question" )
    codec   = command.New( plugin.Flag, "codec",   nil,     "codec method" )
    client  = command.New( plugin.Flag, "client",  nil,     "client method" )
)

var commands = []command.Item{

    { group,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { types,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { codec,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { client,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

}

func (r *stdin) Init() {
    if group := group.GetValue(); group != nil {
        fmt.Println("groupppppppppppppp", group.GetString())
    }

    /*
    registry := paths.Resolve(paths.Data, "registry")
    fmt.Println("hahahahha",registry)
    */

    //利用group codec等,进行初始化

    //init group
    //gValue := group.GetValue()
    //fmt.Println(gValue.GetString())

    // init type
    //tValue := types.GetValue()
    //fmt.Println(tValue.GetString())
/*
    if p := paths.GetArray(); p != nil {
        // init paths
        fmt.Println(p)
    }

    if c := codec.GetMap(); c != nil {
        // init codec
        fmt.Println(c)
        codec.Clear()
    } else {
        // default init
    }

    if c := client.GetMap(); c != nil {
        fmt.Println(c)
        client.Clear()
    } else {
        // default init
    }
*/
    return
}

func (r *stdin) Main() {
    fmt.Println("Start agent file module ...")
    // 编写主要业务逻辑

    c, err := collector.New(r.Log)
    if err != nil {
        //r.Print("Error in initing collector: %s", err)
        fmt.Errorf("Error in initing collector: %s", err)
        r.Exit(-1)
    }

    c.Start()

    return
}

func (r *stdin) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Agents, Name, commands, New)
}
