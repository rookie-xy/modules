package file

import (
    "fmt"
    "time"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/state"
    "github.com/rookie-xy/hubble/plugin"

    "github.com/rookie-xy/modules/agents/file/finder"
    "github.com/rookie-xy/modules/agents/file/configure"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/adapter"
)

const Name  = "file"

type file struct {
    finder     *finder.Finder
    frequency   time.Duration
    log         log.Log
    done        chan struct{}
}

func New(log log.Log) module.Template {
    return &file{
        log: log,
    }
}

var (
    frequency = command.New( module.Flag, "frequency",  3 * time.Second,  "scan frequency method" )
    group     = command.New( module.Flag, "group",     "nginx",     "This option use to group" )
    Type      = command.New( module.Flag, "type",      "log",       "file type, this is use to find some question" )
    paths     = command.New( module.Flag, "paths",     nil,         "File path, its is manny option" )
    excludes  = command.New( module.Flag, "excludes",  nil,         "File path, its is manny option" )
    limit     = command.New( module.Flag, "limit",      uint64(7),        "text finder limit" )
    codec     = command.New( plugin.Flag, "codec",     nil,         "codec method" )
    client    = command.New( plugin.Flag, "client",    nil,         "client method" )
    input     = command.New( plugin.Flag, "input",     nil,         "input method" )
    output    = command.New( plugin.Flag, "output",    nil,         "output method" )
    sincedb   = command.New( plugin.Flag, "sincedb",   nil,         "sincedb method" )
)

var commands = []command.Item{

    { frequency,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { group,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { Type,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { paths,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { excludes,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { limit,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { codec,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { client,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { input,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

     { output,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { sincedb,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },
}

func (f *file) Init() {
    group, Type := group.GetValue(), Type.GetValue()
    if group == nil || Type == nil {
        return
    }

    limit := limit.GetValue()
    if limit == nil {
        return
    }

    key := codec.GetFlag() + "." + codec.GetKey()
    codec, err := factory.Codec(key, f.log, codec.GetValue())
    if err != nil {
        return
    }


	key = input.GetFlag() + "." + input.GetKey()
    input, err := factory.Input(key, f.log, input.GetValue())
    if err != nil {
        return
    }

    configure := configure.Configure{
    	Group:    group.GetString(),
    	Type:     Type.GetString(),
        Paths:    paths.GetValue(),
        Excludes: excludes.GetValue(),
        Limit:    limit.GetUint64(),
        Input:    input,
        Codec:    codec,
        Client:   true,
    }

    if value := client.GetValue(); value != nil {
        key = client.GetFlag() + "." + client.GetKey()
        configure.Output, err = factory.Client(key, f.log, value)
        if err != nil {
            return
        }

    } else {
    	configure.Client = false

        key = output.GetFlag() + "." + output.GetKey()
        configure.Output, err = factory.Output(key, f.log, output.GetValue())
        if err != nil {
            return
        }
    }

    if value := frequency.GetValue(); value != nil {
        f.frequency = value.GetDuration()
    }

 	key = sincedb.GetFlag() + "." + sincedb.GetKey()
    client, err := factory.Forward(key)
    if err != nil {
        return
    }

    finder := finder.New(f.log)
    if err := finder.Init(&configure, adapter.FileSinceDB(client)); err != nil {
        fmt.Println(err)
        return
    }

    f.finder = finder

    return
}

func (f *file) Main() {
    fmt.Println("Start agent file module ...")
    // 编写主要业务逻辑

    //f.wg.Add(1)
    //r.Print("Starting finder of type: %v; id: %v ", p.config.Type, p.ID())

    //onceWg := sync.WaitGroup{}
    //onceWg.Add(1)

    // Add waitgroup to make sure prospectors finished
    run := func(finder *finder.Finder) error {
        defer func() {
            //onceWg.Done()
            close(f.done)
            //f.wg.Wait()
            //f.wg.Done()
        }()

        finder.Find()

        for {
            select {

            case <-f.done:
                //r.Print("Finder ticker stopped")
                return nil
            case <-time.After(f.frequency):
                //r.Debug("finder", "Run finder")
                finder.Find()
            }
        }

        return nil
    }

    run(f.finder)

    return
}

func (f *file) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Agents, Name, commands, New)
}
