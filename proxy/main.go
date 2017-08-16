package outputs

import (
    "fmt"

    "github.com/rookie-xy/hubble/src/command"
    "github.com/rookie-xy/hubble/src/module"
    "github.com/rookie-xy/hubble/src/log"
    "github.com/rookie-xy/hubble/src/register"
    "github.com/rookie-xy/hubble/src/prototype"
    "github.com/rookie-xy/hubble/src/state"

  _ "github.com/rookie-xy/modules/proxy/src/forward"
  _ "github.com/rookie-xy/modules/proxy/src/reverse"
)

const Name = module.Proxy

var (
    proxy = command.New("", Name, nil, "outputs may be many")
)

var commands = []command.Item{

    { proxy,
      command.FILE,
      Name,
      nil,
      state.Enable,
      0,
      nil },

}

type Proxy struct {
    log.Log
    event    chan int
    children []module.Template
}

func New(log log.Log) module.Template {
    new := &Proxy{
        Log: log,
        event: make(chan int, 1),
    }

    register.Observer(Name, new)

    return new
}

func (r *Proxy) Init() {
    // 等待配置更新完成的信号
    <-r.event
    fmt.Println("proxy init")
    //fmt.Println(proxy.Value)
    //return

    if value := proxy.GetArray(); value != nil {
        // key为各个模块名字，value为各个模块配置
        for _, configure := range value {
            // 渲染模块命令
            for name, value := range configure.(map[interface{}]interface{}) {
                // 渲染指令
                if value != nil {
                    for k, v := range value.(map[interface{}]interface{}) {
                        if status := command.File(Name, k.(string), v); status != state.Ok {
                            fmt.Println("command file error", status)
                            //exit(status)
                        }
                    }
                }

                // 安装模块
                key := Name + "." + name.(string)
                module := module.Setup(key, r.Log)
                if module != nil {
                    module.Init()

                } else {
                    fmt.Println("output setup module error")
                    return
                }

                r.Load(module)
            }
        }

    } else {
        fmt.Println("output value is nil")
    }

    return
}

func (r *Proxy) Main() {
    fmt.Println("Start proxy modules ...")
    if len(r.children) < 1 {
        return
    }

    for _, child := range r.children {
        child.Main()
    }

    for ;; {
        select {}
    }
}

func (r *Proxy) Exit(code int) {
    return
}

func (r *Proxy) Update(configure prototype.Object) int {
    if configure == nil {
        return state.Error
    }

    exist := true
    proxy.Value, exist = configure.(map[interface{}]interface{})[Name]
    if !exist {
        fmt.Println("Not found proxy configure")
        return state.Error
    }

    r.event <- 1
    return state.Ok
}

func (r *Proxy) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
