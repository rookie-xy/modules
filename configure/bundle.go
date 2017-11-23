package configure

import (
    "fmt"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
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
    "time"
)

const Name  = module.Configure

var (
    mode  = command.New( "-m", "mode",  "local", "Specifies how to obtain the configuration source" )
    style = command.New( "-s", "style", "yaml", "Specifies the format of the configuration source" )
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

    reload     bool
    observers  []observer.Observer
    event      chan types.Object
    children   []module.Template
}

func New(log log.Log) module.Template {
    new := &Configure{
        Log: log,
        event: make(chan types.Object, 1),
        reload: false,
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

    if obslen := len(r.observers); o != nil && obslen > 0 {
    	fmt.Println("configure modulesssssssssssssssssssssssssss", obslen)

        var mains []func()

        for _, configure := range r.observers {
            if configure.Update(o) != nil {
                break
            }

            if r.reload {
                if module := adapter.ToModuleObserver(configure); module != nil {
                    //fmt.Println("Start exit all the module ... ...")
                   	//time.Sleep(10 * time.Second)
                   	module.Exit(0)

                   	time.Sleep(10 * time.Second)
                   	fmt.Println("All the module is exit, Start init all the module ... ...")
                   	time.Sleep(2 * time.Second)

                    module.Init()

                    //mains = append(mains, module.Main)
                }
            }
        }

        if length := len(mains); length > 0 {
        	fmt.Println("All the module init is finish, Start running all the module ... ...")
        	time.Sleep(7 * time.Second)
        	for _, main := range mains {
                go main()
			}

			r.reload = false
        }
    }

    return
}

func (r *Configure) Update(o types.Object) error {
    _, data, err := r.ValueDecode(o.([]byte), true)
    if err != nil {
        return fmt.Errorf("error", data)
    }

    r.event <- data
    return nil
}

func (r *Configure) Reload(o types.Object) error {
    _, data, err := r.ValueDecode(o.([]byte), true)
    if err != nil {
        return fmt.Errorf("error", data)
    }

    r.reload = true
    r.event <- data
    return nil
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

// 两件事情：
// 1. 解析配置，
// 2. 加载配置，如果是第一次启动则加载配置，如果是跟新配置，则reload所有组件
func (r *Configure) Main() {
    // 启动各个子模块组件
    for _, child := range r.children {
        child.Init()
        go child.Main()
    }

    for ;; {
        select {
        case e := <- r.event:
            // 先通知在reload
            fmt.Println("configure notifyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy")
            r.Notify(e)
        }
    }
}

func (r *Configure) Exit(code int) {
	if num := len(r.children); num > 0 {
        for _, module := range r.children {
            module.Exit(code)
        }
    }
}

func (r *Configure) Load(m module.Template) {
    if m != nil {
        r.children = append(r.children, m)
    }
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}


/*
func (r *Configure) Notify(o types.Object) {
    if o != nil {
        r.update(o)
    }

    if r.reload {

        //reload()
    }
    return
}

func (r *Configure) update(o types.Object) {
    for _, observer := range r.observers {
        if observer.Update(o) != nil {
            break
        }
    }
}
*/

