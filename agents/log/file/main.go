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

    "github.com/rookie-xy/modules/agents/log/collector"
    "github.com/rookie-xy/modules/agents/log/file/scanner"
)

const Name  = "file"

type file struct {
    log        log.Log
    scanner   *scanner.Scanner

    done       chan struct{}
    wg        *sync.WaitGroup
    id         uint64

    frequency  time.Duration
}

func New(log log.Log) module.Template {
    return &file{
        log: log,
    }
}

var (
    group     = command.New( module.Flag, "group",     "nginx", "This option use to group" )
    Type      = command.New( module.Flag, "type",      "log",   "file type, this is use to find some question" )
    paths     = command.New( module.Flag, "paths",     nil,     "File path, its is manny option" )
    excludes  = command.New( module.Flag, "paths",     nil,     "File path, its is manny option" )
    codec     = command.New( plugin.Flag, "codec",     nil,     "codec method" )
    client    = command.New( plugin.Flag, "client",    nil,     "client method" )
    frequency = command.New( module.Flag, "frequency", 10 * time.Second, "scan frequency method" )
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

}

func (r *file) Init() {
    group, Type := group.GetValue(), Type.GetValue()
    if group == nil || Type == nil {
        return
    }

    collector := collector.New(r.log)
    if err := collector.Init(group.GetString(), Type.GetString(),
                             codec, client); err != nil {
        fmt.Println(err)
        return
    }

    scanner := scanner.New(r.log)
    if err := scanner.Init(Name, paths, excludes, collector); err != nil {
        fmt.Println(err)
        return
    }

    r.scanner = scanner

    return
}

func (r *file) Main() {
    fmt.Println("Start agent file module ...")
    // 编写主要业务逻辑

    r.wg.Add(1)
    //r.Print("Starting prospector of type: %v; id: %v ", p.config.Type, p.ID())

    onceWg := sync.WaitGroup{}
    onceWg.Add(1)

    // Add waitgroup to make sure prospectors finished
    run := func() int {
        defer func() {
            onceWg.Done()
            r.stop()
            r.wg.Done()
        }()

        r.Run()

        return state.Ok
    }

    run()

    return
}

func (r *file) Run() {
    // Initial tracker run
    r.scanner.Run()

    for {
        select {

        case <-r.done:
            //r.Print("Collector ticker stopped")
            return

        case <-time.After(r.frequency):
            //r.Debug("collector", "Run collector")
            r.scanner.Run()
        }
    }
}

func (r *file) stop() {
    //r.Print("Stopping Collector: %v", r.ID())
    r.scanner.Stop()
}

func (r *file) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Agents, Name, commands, New)
}
