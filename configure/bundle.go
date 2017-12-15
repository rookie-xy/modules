package configure

import (
    "log"
    "os"

    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/plugin"
    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/memento"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/observer"
    "github.com/rookie-xy/hubble/register"

  l "github.com/rookie-xy/hubble/log"
  . "github.com/rookie-xy/hubble/log/level"

  _ "github.com/rookie-xy/modules/configure/local"
  _ "github.com/rookie-xy/modules/configure/remote"
    "sync"
)

var (
    mode  = command.New("-m", "mode",  "local", "Specifies how to obtain the configuration source" )
    style = command.New("-s", "style", "yaml", "Specifies the format of the configuration source" )
    debug = command.New("-d", "debug", false,       "output detail info")
    level = command.New("-l", "level",   "info",     "output detail info")
    title = command.New("-t", "title",  "[hubble] ", "output detail info")
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
    level  Level

    codec      codec.Codec
    done       chan struct{}
    wg         sync.WaitGroup
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
	        "[" + prefix.GetString() + "] ",
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

    new := &Configure{
        Log: lg,
        level: value,
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

    r.log(ERROR,"configure attach error")
}

func (r *Configure) Notify(o types.Object) {

    if n := len(r.observers); o != nil && n > 0 {

    	var inits, mains []func()
        for _, configure := range r.observers {
            if err := configure.Update(o); err != nil {
                r.log(WARN, err.Error())
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
        if r.reload {
            r.log(DEBUG, "Reload, All core components should be shut down")
        }

        r.reload = false

        if l := len(inits); l > 0 {
        	r.log(DEBUG,"Reload, Initialization all core components")

        	for _, init := range inits {
        	    init()
			}
        }

        if l := len(mains); l > 0 {
        	r.log(DEBUG,"Reload, Run all core components")

        	for _, main := range mains {
                go main()
			}
        }
    }
}

func (r *Configure) Update(o types.Object) error {
    data, err := r.codec.Decode(o.([]byte))
    if err != nil {
        return err
    }

    r.event <- data
    r.log(INFO,"Configure update, ready to load all components")
    return nil
}

func (r *Configure) Reload(o types.Object) error {
	data, err := r.codec.Decode(o.([]byte))
    if err != nil {
        return err
    }

    r.reload = true
    r.event <- data
    r.log(INFO,"Configure reload, ready to exit and load all components")
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
    if domain, ok := plugin.Domain(codec.Name, value.GetString()); ok {
        var err error
        if r.codec, err = factory.Codec(domain, r.Log, nil); err != nil {
            r.log(ERROR, err.Error())
        }
    } else {
        r.log(ERROR, "plugin domain failure")
    }
}

// two things:
// 1. parsing configuration
// 2. load configuration, if it is the first time to start loading
//    the configuration, if it is with the new configuration, then all
//    the components of reload
func (r *Configure) Main() {
    r.log(DEBUG,"Run components for configure")

    defer func() {
        r.wg.Wait()
    }()

    r.wg.Add(len(r.children))
    for _, child := range r.children {
    	if child != nil {
    	    child.Init()
            go func(main func()) {
        	    defer r.wg.Done()
        	    main()

            }(child.Main)

            continue
        }

        r.log(WARN, "Configure child component is nil [main stage]")
    }

    r.log(DEBUG, "Configure all components have started running")

    for {
        select {
        case e := <- r.event:
            r.Notify(e)
        case <- r.done:
            r.log(INFO, "Configure main process is exit")
            return
        }
    }
}

func (r *Configure) Exit(code int) {
    defer func(children []module.Template, exit bool, code int) {
        if exit {
            for _, child := range children {
                if child != nil {
                    child.Exit(code)
                    continue
                }

                r.log(WARN, "Configure child component is nil [exit stage]")
            }

            r.log(DEBUG, "Configure all components have exit")
        }

        r.log(WARN, "Configure no components need to quit")
    } (r.children, len(r.children) > 0, code)

    r.log(INFO,"Exit components for configure")
    close(r.done)
}

func (r *Configure) Load(m module.Template) {
    if m != nil {
        r.children = append(r.children, m)
    }
}

func (r *Configure) log(ll Level, f string, args ...interface{}) {
    l.Print(r.Log, r.level, ll, f, args...)
}

func init() {
    register.Component(module.Worker, Name, commands, New)
}
