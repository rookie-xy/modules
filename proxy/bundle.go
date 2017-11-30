package proxy

import (
    "fmt"
    "sync"

    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
  . "github.com/rookie-xy/hubble/log/level"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/configure"
    "github.com/rookie-xy/hubble/types/value"

  _ "github.com/rookie-xy/modules/proxy/forward"
  _ "github.com/rookie-xy/modules/proxy/sinceDB"
    "github.com/rookie-xy/hubble/errors"
)

var proxy = command.New("", Name, nil, "outputs may be many")

var commands = []command.Item{

    { proxy,
      command.FILE,
      Name,
      Name,
      nil,
      nil },

}

type Proxy struct {
    log.Log
    level     Level

    wg        sync.WaitGroup
    event     chan int
    children  []module.Template
    done      chan struct{}
}

func New(log log.Log) module.Template {
    new := &Proxy{
        Log:   log,
        level: adapter.ToLevelLog(log).Get(),
        event: make(chan int, 1),
    }

    register.Observer(Name, new)
    return new
}

func (r *Proxy) Init() {
    r.log(DEBUG, "Proxy, waiting for a signal to configure the update")

    <-r.event
    r.log(DEBUG,"Initialization components for proxy")

    r.done = make(chan struct{})
    r.children = []module.Template{}

    build := func(scope string, i types.Iterator, load module.Load) error {
        for iterator := i; iterator.Has(); iterator.Next() {
            iterm := iterator.Iterm()
            name := iterm.Key.GetString()


            if v := iterm.Value; v != nil {
                value := value.New(v)
                it := value.GetIterator(nil)
                if it == nil {
                    continue
                }

                for iterator := it; iterator.Has(); iterator.Next() {
                    if iterm := iterator.Iterm(); iterm != nil {
                        key := iterm.Key.GetString()
                        if err := command.File(scope, name, key, iterm.Value); err != nil {
                            return err
                        }
                    }
                }
            }

            namespace := scope + "." + name
            module := module.Setup(namespace, r.Log)
            if module != nil {
                module.Init()

            } else {
                return fmt.Errorf("[%s] module setup error\n", name)
            }

            load(module)
        }

        return nil
    }

    if proxy := proxy.GetValue(); proxy != nil {
        iterator := proxy.GetIterator(nil)
        if iterator != nil {
            if err := build(Name, iterator, r.Load); err != nil {
            	r.log(ERROR,"Initialization components for proxy, " +
            	      "proxy build error %s", err)
                return
            }

            configure.Build = build
        }
    }

    r.log(DEBUG,"Proxy all component initialization completed")
}

func (r *Proxy) Main() {
    r.log(DEBUG,"Run components for proxy")

    defer func() {
        r.wg.Wait()
        close(r.done)
    }()

    r.wg.Add(len(r.children))

    for _, child := range r.children {
        if child != nil {
        	go func(main func()) {
        	    defer r.wg.Done()

        	    main()

            }(child.Main)

            continue
        }

        r.log(WARN, "Proxy child component is nil [main stage]")
    }

    r.log(DEBUG, "Proxy all components have started running")
}

func (r *Proxy) Exit(code int) {
    defer func() {
        <-r.done
        r.log(DEBUG,"Proxy all components have exit")
    }()

    r.log(DEBUG,"Exit components for proxy")

    if n := len(r.children); n > 0 {
        for _, child := range r.children {
        	if child != nil {
                child.Exit(code)
                continue
            }

            r.log(WARN, "Proxy child component is nil [exit stage]")
        }
    }
}

func (r *Proxy) Update(o types.Object) error {
    v := value.New(o)
    if v.GetType() != types.MAP {
        return errors.ErrType
    }

    if value := v.GetMap(); value != nil {
        val, exist := value[Name]
        if !exist {
            return errors.ErrConfigure
        }

        proxy.SetValue(val)
    }

    r.log(DEBUG,"Update proxys configure successful")
    r.event <- 1

    return nil
}

func (r *Proxy) Load(m module.Template) {
    r.children = append(r.children, m)
}

func (r *Proxy) log(l Level, f string, args ...interface{}) {
    log.Print(r.Log, r.level, l, f, args...)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
