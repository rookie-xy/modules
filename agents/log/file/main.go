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

    "github.com/rookie-xy/modules/agents/log/text"
    "github.com/rookie-xy/modules/agents/log/file/scanner"
)

const Name  = "file"

type file struct {
    scanner    *scanner.Scanner
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
    excludes  = command.New( module.Flag, "paths",     nil,              "File path, its is manny option" )
    codec     = command.New( plugin.Flag, "codec",     nil,              "codec method" )
    client    = command.New( plugin.Flag, "client",    nil,              "client method" )
    frequency = command.New( module.Flag, "frequency", 10 * time.Second, "scan frequency method" )
    limit     = command.New( module.Flag, "limit",     7,                "text scanner limit" )
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
    group, Type := group.GetValue(), Type.GetValue()
    if group == nil || Type == nil {
        return
    }

    text := text.New(f.log)
    if err := text.Init(group.GetString(), Type.GetString(),
                                              codec, client); err != nil {
        fmt.Println(err)
        return
    }

    if value := frequency.GetValue(); value != nil {
        f.frequency = value.GetUint64()
    }

    limit := limit.GetValue()
    if limit == nil {
        return
    }

    fscanner := scanner.New(f.log)
    if err := fscanner.Init(Name, paths, excludes,
                                         text, limit.GetUint64()); err != nil {
        fmt.Println(err)
        return
    }

    f.scanner = fscanner

    return
}

func (f *file) Main() {
    fmt.Println("Start agent file module ...")
    // 编写主要业务逻辑

    f.wg.Add(1)
    //r.Print("Starting prospector of type: %v; id: %v ", p.config.Type, p.ID())

    onceWg := sync.WaitGroup{}
    onceWg.Add(1)

    // Add waitgroup to make sure prospectors finished
    run := func(file *scanner.Scanner) int {
        defer func() {
            onceWg.Done()
            close(f.done)
            f.wg.Wait()
            f.wg.Done()
        }()

        // Initial scan run
        file.Scan()

        for {
            select {

            case <-f.done:
                //r.Print("Collector ticker stopped")
                return

            case <-time.After(f.frequency):
                //r.Debug("collector", "Run collector")
                file.Scan()
            }
        }

        return state.Ok
    }

    run(f.scanner)

    return
}

func (f *file) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Agents, Name, commands, New)
}
