package configure

import (
    "unsafe"
    "fmt"

    "github.com/rookie-xy/plugins"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/configure"
    "github.com/rookie-xy/worker/src/observer"
    "github.com/rookie-xy/worker/src/prototype"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/plugin"

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
    plugin.Codec
    observers []observer.Observer
    data prototype.Object
    children []module.Template
}

func New(log log.Log) *Configure {
    return &Configure{
        Log: log,
    }
}

func (r *Configure) Attach(obs observer.Observer) {
    r.observers = append(r.observers, obs)
}

func (r *Configure) Notify() {
    for _, observer := range r.observers {
         observer.Update(r.data)
    }
}

func (r *Configure) Init() {
    if v := config.Value; v != nil {
        key := Name + "_" + v.(string)

        if module := module.Setup(key, r.Log); module != nil {
            r.Load(module)
        }
    }

    if v := format.Value; v != nil {
        if codec := plugins.Codec(v.(string)); codec != nil {
            r.Codec = codec
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

    // 渲染所有配置指令
    for ;; {

        select {

        case e := <-configure.Event:
            fmt.Println(e)
            /*
            data, err := r.Codec.Decode(e)
            if err != nil {
                fmt.Println(data)
            }
            */
            // TODO 解析配置，通知加载三大模块
            // TODO 监听外部启停指令
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
    register.Module(module.Worker, Name, commands, nil)
}
