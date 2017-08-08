package memory

import (
    "fmt"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/channel"
    "github.com/rookie-xy/worker/src/state"

  . "github.com/rookie-xy/modules/channels/src/memory/src"
    "github.com/rookie-xy/modules/channels/src/memory/src/subject"
)

const Name  = "memory"

type memory struct{
    log.Log
    channel.Pull
    subject  *subject.Subject
    filter   *Filter
}

var (
    name   = command.Metas( "", "name", "nginx",    "This option use to group" )
    mode   = command.Metas( "", "mode", "pipeline", "This option use to group" )
    size   = command.Metas( "", "size", "16384",    "file type, this is use to find some question" )
    filter = command.Metas( "", "filter", nil,      "file type, this is use to find some question" )
)

var commands = []command.Item{

    { name,
      command.FILE,
      module.Channels,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { mode,
      command.FILE,
      module.Channels,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { size,
      command.FILE,
      module.Channels,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { filter,
      command.FILE,
      module.Channels,
      command.SetObject,
      state.Enable,
      0,
      nil },

}

func New(log log.Log) module.Template {
    return &memory{
        Log: log,
    }
}

func (r *memory) Init() {

    fmt.Println("wwwwwwwwwwwwwwwwwwwwwww", name.Value, mode.Value, size.Value, filter.Value)

    names := ""
    if v := name.Value; v != nil {
        names = v.(string)
    }
/*
    r.Pull, r.subject = factory.Pull(name), subject.New(r.Log)
    if r.Pull == nil || r.subject == nil {
        fmt.Println("pull or subject init error")
        return
    }
    */

    r.subject = subject.New(r.Log)
    if r.subject == nil {
        fmt.Println("pull or subject init error")
        return
    }

    register.Subject(names, r.subject)

    return
}

func (r *memory) Main() {
    if r.Pull == nil {
        return
    }

    for {
        select {
         /*
        case event := r.Pull.Pull():
            Process(event, r.filter, r.subject)
            */
        }

    }
}

func (r *memory) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Channels, Name, commands, New)
}
