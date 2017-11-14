package file

import (
    "fmt"
    "time"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/plugin"

    "github.com/rookie-xy/modules/agents/file/finder"
    "github.com/rookie-xy/modules/agents/file/configure"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/adapter"
    Output "github.com/rookie-xy/hubble/output"
    "github.com/rookie-xy/hubble/types/value"
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
    Type      = command.New( module.Flag, "type",      "log",       "source type, this is use to find some question" )
    paths     = command.New( module.Flag, "paths",     nil,         "File path, its is manny option" )
    excludes  = command.New( module.Flag, "excludes",  nil,         "File path, its is manny option" )
    limit     = command.New( module.Flag, "limit",      uint64(7),        "text finder limit" )
    codec     = command.New( plugin.Flag, "codec",     nil,         "codec method" )
    client    = command.New( plugin.Flag, "client",    nil,         "client method" )
    input     = command.New( plugin.Flag, "input",     nil,         "input method" )
    output    = command.New( plugin.Flag, "output",    nil,         "output method" )
    sinceDB   = command.New( plugin.Flag, "client.sinceDB",   nil,         "sinceDB method" )
)

var commands = []command.Item{

    { frequency,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

    { group,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

    { Type,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

    { paths,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

    { excludes,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

    { limit,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

    { codec,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

    { client,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

    { input,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

     { output,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
      nil },

    { sinceDB,
      command.FILE,
      module.Agents,
      Name,
      command.SetObject,
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
    	fmt.Println("agent source codec: ", err)
        return
    }

    pluginName := input.GetFlag() + "." + input.GetKey()
    input, err := factory.Input(pluginName, f.log, input.GetValue())
    if err != nil {
    	fmt.Println("agent source input: ", err)
        return
    }

    configure := configure.Configure{
    	Group:    group.GetString(),
    	Type:     Type.GetString(),
        Paths:    paths.GetValue(),
        Excludes: excludes.GetValue(),
        Input:    input,
        Codec:    codec,
    }

    if limit, err := limit.GetUint64(); err != nil {
        fmt.Println(err)
        return
    } else {
        configure.Limit = limit
    }

    if value := client.GetValue(); value != nil {
   	    configure.Client = command.New(
            client.GetFlag(),
            client.GetKey(),
            client.GetObject(),
           "")

        configure.SinceDB = command.New(
            sinceDB.GetFlag(),
            sinceDB.GetKey(),
            sinceDB.GetObject(),
           "")
    } else {
   	    configure.Output = command.New(
            output.GetFlag(),
            output.GetKey(),
            output.GetObject(),
           "")
    }

    if value := frequency.GetValue(); value != nil {
        if duration, err := value.GetDuration(); err != nil {
        	fmt.Println(err)
        	return
        } else {
            f.frequency = duration
        }
    }

 	key = sinceDB.GetFlag() + "." + Output.Name + "." + "sinceDB"
    sinceDB, err := factory.Output(key, f.log, value.New(sinceDB.GetKey()))
    if err != nil {
        fmt.Println("agent file sinceDB: ", err)
        return
    }

    finder := finder.New(f.log)
    if err := finder.Init(&configure, adapter.FileSinceDB(sinceDB)); err != nil {
        fmt.Println("agent file finder init: ", err)
        return
    }

    f.finder = finder

    return
}

func (f *file) Main() {
    fmt.Println("Start agent file module ...")
    // 编写主要业务逻辑

    run := func(finder *finder.Finder) error {
        defer func() {
            //close(f.done)
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
