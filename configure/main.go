package configure

import (
    "fmt"

    "github.com/rookie-xy/hubble/src/command"
    "github.com/rookie-xy/hubble/src/module"
    "github.com/rookie-xy/hubble/src/configure"
    "github.com/rookie-xy/hubble/src/observer"
    "github.com/rookie-xy/hubble/src/prototype"
    "github.com/rookie-xy/hubble/src/register"
    "github.com/rookie-xy/hubble/src/log"
    "github.com/rookie-xy/hubble/src/codec"
    "github.com/rookie-xy/hubble/src/state"
    "github.com/rookie-xy/hubble/src/factory"

  _ "github.com/rookie-xy/modules/configure/src/file"
  _ "github.com/rookie-xy/modules/configure/src/zookeeper"

)

const Name  = module.Configure

var (
    config = command.New( "-c", "config", "file", "Specifies how to obtain the configuration file" )
    format = command.New( "-f", "format", "yaml", "Specifies the format of the configuration file" )
)

var commands = []command.Item{

    { config,
      command.LINE,
      module.Configure,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { format,
      command.LINE,
      module.Configure,
      command.SetObject,
      state.Enable,
      0,
      nil },

}

type Configure struct {
    log.Log
    codec.Codec

    observers  []observer.Observer
    data       prototype.Object
    children   []module.Template
}

func New(log log.Log) module.Template {
   new := &Configure{
        Log: log,
    }

    register.Subject(Name, new)

    return new
}

func (r *Configure) Attach(o observer.Observer) {
    if o != nil {
        r.observers = append(r.observers, o)
        return
    }

    fmt.Println("attach error")
    return
}

func (r *Configure) Notify() {
    if r.data == nil {
        return
    }

    //fmt.Println(r.data)
    r.update()
}

func (r *Configure) update() {
    for _, observer := range r.observers {
        if observer.Update(r.data) == state.Error {
            break
        }
    }
}

func (r *Configure) Init() {
    key := Name + "." + config.GetString()
    if module := module.Setup(key, r.Log); module != nil {
        r.Load(module)
    }

    if codec, err := factory.Codec(format.GetString(), r.Log, nil); err != nil {
        fmt.Println(err)
        return

    } else {
        r.Codec = codec
    }

    return
}

func (r *Configure) Main() {
    // 启动各个子模块组件
    for _, child := range r.children {
        child.Init()
        go child.Main()
    }

    for ;; {

        select {

        case e := <-configure.Event:
            var err error

            r.data, err = r.Decode(e)
            if err != nil {
                fmt.Println("error", r.data)
                return
            }

            r.Notify()

        default:

        }
    }

    return
}

func (r *Configure) Exit(code int) {
    for _, module := range r.children {
        module.Exit(code)
    }

    //r.cycle.Quit()
    return
}

func (r *Configure) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
