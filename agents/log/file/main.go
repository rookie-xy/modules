package file

import (
    "fmt"
    "sync"
    "time"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/state"
    "github.com/rookie-xy/hubble/plugin"

//    "github.com/rookie-xy/modules/agents/log/collector"
    "github.com/rookie-xy/modules/agents/log/file/finder"
)

const Name  = "file"

type file struct {
    finder     *finder.Finder
    frequency   time.Duration
    log         log.Log

    done        chan struct{}
    wg         *sync.WaitGroup
}

func New(log log.Log) module.Template {
    return &file{
        log: log,
    }
}

var (
    group     = command.New( module.Flag, "group",     "nginx",          "This option use to group" )
    Type      = command.New( module.Flag, "type",      "log",            "file type, this is use to find some question" )
    paths     = command.New( module.Flag, "paths",     nil,              "File path, its is manny option" )
    excludes  = command.New( module.Flag, "excludes",   nil,              "File path, its is manny option" )
    codec     = command.New( plugin.Flag, "codec",     nil,              "codec method" )
    client    = command.New( plugin.Flag, "client",    nil,              "client method" )
    frequency = command.New( module.Flag, "frequency", 3 * time.Second, "scan frequency method" )
    limit     = command.New( module.Flag, "limit",     uint64(7),                "text finder limit" )
)

var commands = []command.Item{

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

    { frequency,
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

}

func (f *file) Init() {
/*
    group, Type := group.GetValue(), Type.GetValue()
    if group == nil || Type == nil {
        return
    }

    collector := collector.New(f.log)
    if err := collector.Init(group.GetString(), Type.GetString(),
                                              codec, client); err != nil {
        fmt.Println(err)
        return
    }
*/
    if value := frequency.GetValue(); value != nil {
        f.frequency = value.GetDuration()
    }

    limit := limit.GetValue()
    if limit == nil {
        return
    }

    finder := finder.New(f.log)
    if err := finder.Init(Name, paths.GetValue(), excludes.GetValue(),
                                nil, limit.GetUint64()); err != nil {
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

        // Initial finder run
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
