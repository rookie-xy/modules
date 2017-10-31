package configure

import (
    "fmt"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/state"
    "github.com/rookie-xy/hubble/observer"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/log"
//    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/memento"

  _ "github.com/rookie-xy/modules/configure/local"
  _ "github.com/rookie-xy/modules/configure/remote"
    "github.com/rookie-xy/hubble/plugin"
    "github.com/rookie-xy/hubble/codec"
)

const Name  = module.Configure

var (
    mode  = command.New( "-m", "mode",  "local", "Specifies how to obtain the configuration file" )
    style = command.New( "-s", "style", "yaml", "Specifies the format of the configuration file" )
)

var commands = []command.Item{

    { mode,
      command.LINE,
      module.Configure,
      Name,
      command.SetObject,
      nil },

    { style,
      command.LINE,
      module.Configure,
      Name,
      command.SetObject,
      nil },

}

type Configure struct {
    log.Log
    adapter.ValueCodec

    observers  []observer.Observer
    event      chan types.Object
    children   []module.Template
}

func New(log log.Log) module.Template {
    new := &Configure{
        Log: log,
        event: make(chan types.Object, 1),
    }

    register.Subject(Name, new)
    register.Observer(Name, new)

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

func (r *Configure) Notify(o types.Object) {
    if o != nil {
        r.update(o)
    }

    //fmt.Println(r.data)
    return
}

func (r *Configure) update(o types.Object) {
    for _, observer := range r.observers {
        if observer.Update(o) == state.Error {
            break
        }
    }
}

func (r *Configure) Update(o types.Object) int {
    _, data, err := r.ValueDecode(o.([]byte), true)
    if err != nil {
        fmt.Println("error", data)
        return state.Error
    }

    r.event <- data

    return state.Ok
}

func (r *Configure) Init() {
    if value := mode.GetValue(); value != nil {
        memento.Name = Name + "." + value.GetString()
        if module := module.Setup(memento.Name, r.Log); module != nil {
            r.Load(module)
        }
    }

    value := style.GetValue()
    if value == nil {
    	fmt.Println("style get value error")
        return
    }

    pluginName := plugin.Flag + "." + codec.Name + "." + value.GetString()

    if codec, err := factory.Codec(pluginName, r.Log, nil); err != nil {
        fmt.Println(err)
        return

    } else {
        r.ValueCodec = adapter.ToValueCodec(codec)
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

        case e := <- r.event:
            r.Notify(e)

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
