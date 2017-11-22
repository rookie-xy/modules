package local

import (
    "os"
    "fmt"

    "github.com/fsnotify/fsnotify"

    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/module"
    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/register"
    "github.com/rookie-xy/hubble/configure"
)

const Name  = "local"

var (
    file = command.New("-f", "source", "./usr/local/conf/hubble.yaml", "If you want to " +
                                     "get locally, you need to specify the prolocal path")
)

var commands = []command.Item{

    { file,
      command.LINE,
      module.Configure,
      Name,
      command.SetObject,
      nil },

}

type local struct {
    log.Log
   *configure.Configure

    name      string
    size      int64

    watcher  *fsnotify.Watcher
    done      chan struct{}
}

func New(log log.Log) module.Template {
    new := &local{
        Log: log,
        Configure: configure.New(log),
    }

    register.Subject(module.Configure + "." + Name, new.Configure)
    return new
}

func (r *local) Init() {
    // 初始化文件监视器，监控配置文件
    resource := ""

    if value := file.GetValue(); value != nil {
        resource = value.GetString()
    }

    localInfo, err := os.Stat(resource)
    if err != nil {
        if os.IsNotExist(err) {
	    fmt.Println("a local or directory does not exist")

        } else if os.IsPermission(err) {
            fmt.Println("permission is denied")

        } else {
            fmt.Println(err)
        }

        return
    }

    r.name = resource
    r.size = localInfo.Size()

    r.watcher, err = fsnotify.NewWatcher()
    if err != nil {
        fmt.Println(err)
        return

    } else {
        r.watcher.Add(resource)
    }

    return
}

func (r *local) Main() {
    var char []byte
    local, err := os.OpenFile(r.name, os.O_RDWR, 0777)
    if err != nil {
        fmt.Println(err)
    }

    char = make([]byte, r.size)

    if size, err := local.Read(char); err != nil {
        if size != int(r.size) {
//            r.Print("size is not r.size")

        } else {
//            r.Print(err)
        }

        local.Close()
        return
    }

    r.Notify(char)

    local.Close()

    // 发现文件变更，通知给其他模块
    for {
        select {

        case event := <-r.watcher.Events:
            if event.Op & fsnotify.Write == fsnotify.Write {
                fmt.Println("OKKKKKKKKKKKKKKKKMENGSHIIIIIIIIIIIIIIIIIII updateeeeeeeeeeeeeeee")
                //r.Notify(char)
                r.Reload(char)
            }

        case err := <-r.watcher.Errors:
            //r.Print(err)
            fmt.Println(err)

        case <-r.done:
        	return
        }
    }
}

func (r *local) Exit(code int) {
	close(r.done)
}

func init() {
    register.Module(module.Configure, Name, commands, New)
}
