package memory

import (
    "fmt"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
//    "github.com/rookie-xy/worker/src/factory"
    "github.com/rookie-xy/worker/src/channel"

  . "github.com/rookie-xy/modules/channels/src/memory/src"
    "github.com/rookie-xy/modules/channels/src/memory/src/subject"
"github.com/rookie-xy/worker/src/state"
)

const Name  = "memory"

type memoryChannel struct{
    log.Log
    channel.Pull
    subject  *subject.Subject
    filter   *Filter
}

var (
    name   = command.Metas( "", "name", "nginx",    state.On, "This option use to group" )
    mode   = command.Metas( "", "mode", "pipeline", state.On, "This option use to group" )
    size   = command.Metas( "", "size", "16384",    state.On, "file type, this is use to find some question" )
    filter = command.Metas( "", "filter", nil,      state.On, "file type, this is use to find some question" )
)

var commands = []command.Item{

    { name,
      command.FILE,
      module.Channels,
      command.SetObject,
      0,
      nil },

    { mode,
      command.FILE,
      module.Channels,
      command.SetObject,
      0,
      nil },

    { size,
      command.FILE,
      module.Channels,
      command.SetObject,
      0,
      nil },

    { filter,
      command.FILE,
      module.Channels,
      command.SetObject,
      0,
      nil },

}

func New(log log.Log) module.Template {
    return &memoryChannel{
        Log: log,
    }
}

func (r *memoryChannel) Init() {

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

func (r *memoryChannel) Main() {
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

func (r *memoryChannel) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Channels, Name, commands, New)
}
