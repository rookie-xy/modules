package file

import (
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
  .  "github.com/rookie-xy/hubble/log/level"
)

const Name  = "file"

type file struct {
    configure  *configure.Configure
    finder     *finder.Finder
    frequency   time.Duration
    log.Log
    level       Level
    done        chan struct{}
    quit        chan struct{}
}

func New(log log.Log) module.Template {
    return &file{
        Log: log,
        level: adapter.ToLevelLog(log).Get(),
        done: make(chan struct{}),
        quit: make(chan struct{}),
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
    f.log(DEBUG,"Initialization %s file component for agent\n", group.GetValue().GetString())

    group, Type := group.GetValue(), Type.GetValue()
    if group == nil || Type == nil {
        return
    }

    limit := limit.GetValue()
    if limit == nil {
        return
    }

    key := codec.GetFlag() + "." + codec.GetKey()
    decoder, err := factory.Decoder(key, f.Log, codec.GetValue())
    if err != nil {
        f.log(ERROR, Name +"; codec: %s", err)
        return
    }

    pluginName := input.GetFlag() + "." + input.GetKey()
    input, err := factory.Input(pluginName, f.Log, input.GetValue())
    if err != nil {
    	f.log(ERROR, Name +"; input: %s", err)
        return
    }

    configure := configure.Configure{
    	Group:    group.GetString(),
    	Type:     Type.GetString(),
        Paths:    paths.GetValue(),
        Excludes: excludes.GetValue(),
        Output:   true,
    }

    if limit, err := limit.GetUint64(); err != nil {
        f.log(ERROR, Name +"; limit: %s", err)
        return
    } else {
        configure.Limit = limit
    }

    if key, ok := plugin.Name(output.GetKey()); ok {
        configure.Client, err = factory.Output(key, f.Log, output.GetValue())
        if err != nil {
            f.log(ERROR, Name + "; output: %s", err)
            return
        }
    }

    if value := client.GetValue(); value != nil {
        if key, ok := plugin.Name(client.GetKey()); ok {
            configure.Client, err = factory.Client(key, f.Log, client.GetValue())
            if err != nil {
                f.log(ERROR, Name+"; client: %s", err)
                return
            }
        }
        configure.Output = false
    }

    if value := frequency.GetValue(); value != nil {
        if duration, err := value.GetDuration(); err != nil {
        	f.log(ERROR, Name +"; duration: %s", err)
        	return
        } else {
            f.frequency = duration
        }
    }

 	key = sinceDB.GetFlag() + "." + Output.Name + "." + "sinceDB"
    sinceDB, err := factory.Output(key, f.Log, value.New(sinceDB.GetKey()))
    if err != nil {
        f.log(ERROR, Name +"; sinceDB: %s", err)
        return
    }

    configure.SinceDB = sinceDB

    finder := finder.New(f.Log)
    if err := finder.Init(input, decoder, &configure, adapter.FileSinceDB(sinceDB), f.log); err != nil {
        f.log(ERROR, Name +"; finder init: %s", err)
        return
    }
    f.finder = finder
    f.configure = &configure

    f.log(DEBUG, Name +"; Agent %s file component initialization completed", configure.Group)
}

func (f *file) Main() {
    f.log(DEBUG, Name +"; Run %s file component for agent\n", f.configure.Group)

    defer func(finder *finder.Finder) {
        finder.Wait()
        close(f.done)
    }(f.finder)

    run := func(finder *finder.Finder) error {
        defer func() {
            finder.Stop()
        }()

        finder.Find()

        f.log(DEBUG, Name +"; Agent %s file component have started running\n", f.configure.Group)
        for {
            select {

            case <-f.quit:
                f.log(INFO, Name +"; Finder ticker stopped")
                return nil
            case <-time.After(f.frequency):
                f.log(DEBUG, Name +"; Run finder")
                finder.Find()
            }
        }

        return nil
    }

    if err := run(f.finder); err != nil {
        f.log(ERROR, Name +"; finder: %s", err)
        return
    }
}

func (f *file) Exit(code int) {
    defer func() {
        <-f.done
        f.log(DEBUG, Name +"; Agent file component have exit")
    }()

    f.log(INFO, Name +"; Exit %s file component for agent", f.configure.Group)
    close(f.quit)
}

func (f *file) log(l Level, fmt string, args ...interface{}) {
    log.Print(f.Log, f.level, l, fmt, args...)
}

func init() {
    register.Component(module.Agents, Name, commands, New)
}
