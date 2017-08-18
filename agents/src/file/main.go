package file

import (
    "fmt"

    "github.com/rookie-xy/hubble/src/command"
    "github.com/rookie-xy/hubble/src/module"
    "github.com/rookie-xy/hubble/src/register"
    "github.com/rookie-xy/hubble/src/log"
    "github.com/rookie-xy/hubble/src/state"
    "github.com/rookie-xy/hubble/src/plugin"
)

const Name  = "file"

type file struct{
    log.Log
}

func New(log log.Log) module.Template {
    return &file{
        Log: log,
    }
}

var (
    group   = command.New( module.Flag, "group",   "nginx", "This option use to group" )
    types   = command.New( module.Flag, "type",    "log",   "file type, this is use to find some question" )
    paths   = command.New( module.Flag, "paths",   nil,     "File path, its is manny option" )
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

    { paths,
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

func (r *file) Init() {
    //利用group codec等,进行初始化
/*
    //init group
    fmt.Println(group.GetString())

    // init type
    fmt.Println(types.GetString())

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

func (r *file) Main() {
    fmt.Println("Start agent file module ...")
    // 编写主要业务逻辑
}

func (r *file) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Agents, Name, commands, New)
}
