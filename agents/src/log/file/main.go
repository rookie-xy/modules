package file

import (
    "fmt"
    "sync"

    "github.com/rookie-xy/hubble/src/command"
    "github.com/rookie-xy/hubble/src/module"
    "github.com/rookie-xy/hubble/src/register"
    "github.com/rookie-xy/hubble/src/log"
    "github.com/rookie-xy/hubble/src/state"
    "github.com/rookie-xy/hubble/src/plugin"

    "github.com/rookie-xy/modules/agents/src/log/api"
    "github.com/rookie-xy/modules/agents/src/log/file/src/collector"
)

const Name  = "file"

type file struct {
    log.Log
   	api.Collector
   	done       chan struct{}
   	wg        *sync.WaitGroup
   	id         uint64
}

func New(log log.Log) module.Template {
    return &file{
        Log: log,
    }
}

var (
    group   = command.New( module.Flag, "group",   "nginx", "This option use to group" )
    types   = command.New( module.Flag, "type",    "log",   "file type, this is use to find some question" )
    paths   = command.New( module.Flag, "paths",   nil,     "File path, its is manny option" )
    codec   = command.New( plugin.Flag, "codec",   nil,     "codec method" )
    client  = command.New( plugin.Flag, "client",  nil,     "client method" )
)

var commands = []command.Item{

    { group,
      command.FILE,
      module.Agents,
      command.SetObject,
      state.Enable,
      0,
      nil },

    { types,
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

}

func (r *file) Init() {
    if group := group.GetValue(); group != nil {
        fmt.Println("groupppppppppppppp", group.GetString())
    }

    r.Collector = collector.New()

    /*
    registry := paths.Resolve(paths.Data, "registry")
    fmt.Println("hahahahha",registry)
    */

    //利用group codec等,进行初始化

    //init group
    //gValue := group.GetValue()
    //fmt.Println(gValue.GetString())

    // init type
    //tValue := types.GetValue()
    //fmt.Println(tValue.GetString())
/*
    if p := paths.GetArray(); p != nil {
        // init paths
        fmt.Println(p)
    }

    if c := codec.GetMap(); c != nil {
        // init codec
        fmt.Println(c)
        codec.Clear()
    } else {
        // default init
    }

    if c := client.GetMap(); c != nil {
        fmt.Println(c)
        client.Clear()
    } else {
        // default init
    }
*/
    return
}

func (r *file) Main() {
    fmt.Println("Start agent file module ...")
    // 编写主要业务逻辑

	   r.wg.Add(1)
	   //r.Print("Starting prospector of type: %v; id: %v ", p.config.Type, p.ID())

	   onceWg := sync.WaitGroup{}
	   if r.Once {
		      // Make sure start is only completed when Run did a complete first scan
		      defer onceWg.Wait()
				}

	   onceWg.Add(1)
	   // Add waitgroup to make sure prospectors finished
	   go func() {
		      defer func() {
			         onceWg.Done()
			         r.stop()
			         r.wg.Done()
		      }()

		     r.Run()
	   }()

    return
}

func (r *file) Exit(code int) {
    // 退出
}

func init() {
    register.Module(module.Agents, Name, commands, New)
}