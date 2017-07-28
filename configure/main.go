package configure

import (
    "unsafe"
    "fmt"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/configure"
    "github.com/rookie-xy/worker/src/observer"
    "github.com/rookie-xy/worker/src/prototype"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/codec"
    "github.com/rookie-xy/worker/src/state"
    "github.com/rookie-xy/worker/src/factory"

  _ "github.com/rookie-xy/modules/configure/src/file"
  _ "github.com/rookie-xy/modules/configure/src/zookeeper"

)

const Name  = module.Configure

var (
    config = &command.Meta{ "-c", "config", "file", "Specifies how to obtain the configuration file" }
    format = &command.Meta{ "-f", "format", "yaml", "Specifies the format of the configuration file" }
)

var commands = []command.Item{

    { config,
      command.LINE,
      module.Configure,
      command.SetObject,
      unsafe.Offsetof(config.Value),
      nil },

    { format,
      command.LINE,
      module.Configure,
      command.SetObject,
      unsafe.Offsetof(format.Value),
      nil },

}

type Configure struct {
    log.Log
    codec     codec.Codec
    observers []observer.Observer
    data      prototype.Object
    children []module.Template
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

    fmt.Println(r.data)
/*
    switch v := r.data.(type) {

    case []interface{}:
        //fmt.Println("A")

    case map[interface{}]interface{}:
        for key, value := range v {
            r.Update()
        }
    }
    */
}

func (r *Configure) Update() {
    for _, observer := range r.observers {
        status := observer.Update("inputs", r.data)
        if status == state.Ok {
            break

        } else if status == state.Declined {
            continue

        } else if status == state.Error {
            break
        }
    }
}

func (r *Configure) Init() {
    if v := config.Value; v != nil {
        key := Name + "." + v.(string)

        if module := module.Setup(key, r.Log); module != nil {
            r.Load(module)
        }
    }

    if v := format.Value; v != nil {

        config := &codec.Config{
            Name: v.(string),
        }

        if codec, err := factory.Codec(config); err != nil {
            fmt.Println(err)
            return

        } else {
            r.codec = codec
        }
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
            r.data, err = r.codec.Decode(e)
            if err != nil {
                fmt.Println("error", r.data)
                return
            }

            fmt.Println(r.data)
            fmt.Println("yuezhanggggggggggggggggggggggggggggg", len(r.observers))

            //r.Notify()

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
