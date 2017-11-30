package configure

import (
    "fmt"
    "log"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/observer"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/types"
  l "github.com/rookie-xy/hubble/log"
  . "github.com/rookie-xy/hubble/log/level"
//    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/memento"

  _ "github.com/rookie-xy/modules/configure/local"
  _ "github.com/rookie-xy/modules/configure/remote"
    "github.com/rookie-xy/hubble/plugin"
    "github.com/rookie-xy/hubble/codec"
    "os"
)

const Name  = module.Configure

var (
    mode  = command.New("-m", "mode",  "local", "Specifies how to obtain the configuration source" )
    style = command.New("-s", "style", "yaml", "Specifies the format of the configuration source" )
    debug = command.New("-d", "debug", false,       "output detail info")
    level = command.New("-l", "level",   "info",     "output detail info")
    title = command.New("-t", "title",  "[hubble]", "output detail info")
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

    { debug,
      command.LINE,
      module.Worker,
      Name,
      command.SetObject,
      nil },

    { level,
      command.LINE,
      module.Worker,
      Name,
      command.SetObject,
      nil },

    { title,
      command.LINE,
      module.Worker,
      Name,
      command.SetObject,
      nil },

}

type Configure struct {
    l.Log
    adapter.ValueCodec

    reload     bool
    observers  []observer.Observer
    event      chan types.Object
    children   []module.Template
}

func New(lg l.Log) module.Template {
	prefix  := title.GetValue()
	verbose := debug.GetValue()
	level   := level.GetValue()

	this := &l.Logger{
		Logger: log.New(
	        os.Stderr,
	        prefix.GetString(),
            log.LstdFlags | log.Lmicroseconds,
        ),
    }
    this.Set(INFO)

    value, err := l.Parse(level.GetString(), verbose.GetBool())
    if err != nil {
        this.Print(ERROR, err.Error())
        return nil
	} else {
        this.Set(value)
    }

    adapter.ToLevelLog(lg).Copy(this)
    lg.Output(3, "heiheiehihhhhhhhhhhhhhhhhhhhhhhhh")
    this.Output(3, "aaaaaaaaaaaaaaa")

    new := &Configure{
        Log: lg,
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

    if l := len(r.observers); o != nil && l > 0 {

    	var inits, mains []func()

        for _, configure := range r.observers {
            if configure.Update(o) != nil {
                break
            }

            if r.reload {
                if module := adapter.ToModuleObserver(configure); module != nil {
                    module.Exit(0)

                    inits = append(inits, module.Init)
                    mains = append(mains, module.Main)
                }
            }
        }

        r.reload = false

        if l := len(inits); l > 0 {
        	fmt.Println("Initialization all core components")
        	for _, init := range inits {
        	    init()
			}
        }

        if l := len(mains); l > 0 {
        	fmt.Println("Run all core components")
        	for _, main := range mains {
                go main()
			}
        }
    }
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
