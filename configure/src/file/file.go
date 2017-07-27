package file

import (
    "os"
    "fmt"
    "unsafe"

    "github.com/fsnotify/fsnotify"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/configure"
)

const Name  = "file"

var (
    path = &command.Meta{ "-p", "path", "./usr/local/conf/worker.yaml", "If you want to " +
                                        "get locally, you need to specify the profile path" }
)

var commands = []command.Item{

    { path,
      command.LINE,
      module.Configure,
      command.SetObject,
      unsafe.Offsetof(path.Value),
      nil },

}

type file struct {
    log.Log

    name      string
    size      int64

    watcher  *fsnotify.Watcher
}

func New(log log.Log) module.Template {
    return &file{
        Log: log,
    }
}

func (r *file) Init() {
    // 初始化文件监视器，监控配置文件
    resource := ""

    if v := path.Value; v != nil {
        resource = v.(string)
    }

    fileInfo, err := os.Stat(resource)
    if err != nil {
        if os.IsNotExist(err) {
	           fmt.Println("a file or directory does not exist")

        } else if os.IsPermission(err) {
            fmt.Println("permission is denied")

        } else {
            fmt.Println(err)
        }

        return
    }

    r.name = resource
    r.size = fileInfo.Size()

    r.watcher, err = fsnotify.NewWatcher()
    if err != nil {
        fmt.Println(err)
        return

    } else {
        r.watcher.Add(resource)
    }

    return
}

func (r *file) Main() {
    var char []byte
    file, err := os.OpenFile(r.name, os.O_RDWR, 0777)
    if err != nil {
        fmt.Println(err)
    //    r.Print("OpenFile error")
    }

    char = make([]byte, r.size)

    if size, err := file.Read(char); err != nil {
        if size != int(r.size) {
//            r.Print("size is not r.size")

        } else {
//            r.Print(err)
        }

        file.Close()
        return
    }

    configure.Event <- char
    file.Close()

    // 发现文件变更，通知给其他模块
    for {
        select {

        case event := <-r.watcher.Events:
            if event.Op & fsnotify.Write == fsnotify.Write {
                configure.Event <- char
            }

        case err := <-r.watcher.Errors:
            r.Print(err)
        }
    }

    return
}

func (r *file) Exit(code int) {
    //r.cycle.Quit()
    return
}

func init() {
    register.Module(module.Configure, Name, commands, New)
}
