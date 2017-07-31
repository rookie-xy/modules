package memory

import (
    "unsafe"
    "fmt"
    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/factory"
    "github.com/rookie-xy/worker/src/channel"

  . "github.com/rookie-xy/modules/channels/src/memory/src"
    "github.com/rookie-xy/modules/channels/src/memory/src/subject"
)

const Name  = "memory"

type memoryChannel struct{
    log.Log
    channel.Pull
    subject  *subject.Subject
    filter   *Filter
}

var (
    name = &command.Meta{ "", "name", "nginx", "This option use to group" }
    size = &command.Meta{ "", "size", "16384", "file type, this is use to find some question" }
)

var commands = []command.Item{

    { name,
      command.FILE,
      module.Channels,
      command.SetObject,
      unsafe.Offsetof(name.Value),
      nil },

    { size,
      command.FILE,
      module.Channels,
      command.SetObject,
      unsafe.Offsetof(size.Value),
      nil },

}

func New(log log.Log) module.Template {
    return &memoryChannel{
        Log: log,
    }
}

func (r *memoryChannel) Init() {
    name = ""
    if v := name.Value; v != nil {
        name = v.(string)
    }

    r.Pull, r.subject = factory.Pull(name), subject.New(r.Log)
    if r.Pull == nil || r.subject == nil {
        fmt.Println("pull or subject init error")
        return
    }

    register.Subject(name, r.subject)

    return
}

func (r *memoryChannel) Main() {
    if r.Pull == nil {
        return
    }

    for {
        select {

        case event := r.Pull():
            Process(event, r.filter, r.subject)
        }
    }
}

func (r *memoryChannel) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Channels, Name, commands, New)
}
