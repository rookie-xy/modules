package inputs

import (
    "fmt"

    "github.com/rookie-xy/hubble/src/command"
    "github.com/rookie-xy/hubble/src/module"
    "github.com/rookie-xy/hubble/src/log"
    "github.com/rookie-xy/hubble/src/prototype"
    "github.com/rookie-xy/hubble/src/register"
    "github.com/rookie-xy/hubble/src/state"

  _ "github.com/rookie-xy/modules/agent/src/file"
)

const Name = module.Inputs

var (
    inputs = command.New("", Name, nil, "inputs may be many")
)

var commands = []command.Item{

    { inputs,
      command.FILE,
      Name,
      nil,
      state.Enable,
      0,
      nil },

}

type Input struct {
    log.Log
    event chan int
    children []module.Template
}

func New(log log.Log) module.Template {
    new := &Input{
        Log: log,
        event: make(chan int, 1),
    }

    register.Observer(Name, new)

    return new
}

func (r *Input) Update(configure prototype.Object) int {
    if configure == nil {
        return state.Error
    }

    exist := true
    inputs.Value, exist = configure.(map[interface{}]interface{})[Name]
    if !exist {
        fmt.Println("Not found inputs configure")
        return state.Error
    }

    r.event <- 1
    return state.Ok
}

func (r *Input) Init() {
    // 等待配置更新完成的信号
    <-r.event
    fmt.Println("input init")
    //fmt.Println(inputs.Value)

    if v := inputs.Value; v != nil {
        // key为各个模块名字，value为各个模块配置
        for _, configure := range v.([]interface{}) {
            // 渲染模块命令
            for name, value := range configure.(map[interface{}]interface{}) {
                // 渲染指令
                for k, v := range value.(map[interface{}]interface{}) {
                    if status := command.File(Name, k.(string), v); status != state.Ok {
                        fmt.Println("command file error", status)
                        //exit(status)
                    }
                }

                // 安装模块
                key := Name + "." + name.(string)
                module := module.Setup(key, r.Log)
                if module != nil {
                    module.Init()

                } else {
                    fmt.Println("inputs setup module error")
                    return
                }

                r.Load(module)
            }
        }

    } else {
        fmt.Println("input value is nil")
    }

    return
}

func (r *Input) Main() {
    // 启动各个组件
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

func (r *Input) Exit(code int) {
    /*
    for _, module := range r.children {
        module.Exit()
    }

    //r.cycle.Quit()
    return
    */
}

func (r *Input) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
