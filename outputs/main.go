package outputs

import (
    "fmt"

    "github.com/rookie-xy/worker/src/command"
    "github.com/rookie-xy/worker/src/module"
    "github.com/rookie-xy/worker/src/log"
    "github.com/rookie-xy/worker/src/register"
    "github.com/rookie-xy/worker/src/prototype"
    "github.com/rookie-xy/worker/src/state"

  _ "github.com/rookie-xy/modules/outputs/src/transport"
)

const Name = module.Outputs

var (
    outputs = command.Metas("", Name, nil, "outputs may be many")
)

var commands = []command.Item{

    { outputs,
      command.FILE,
      Name,
      nil,
      state.Enable,
      0,
      nil },

}

type Output struct {
    log.Log
    event    chan int
    children []module.Template
}

func New(log log.Log) module.Template {
    new := &Output{
        Log: log,
        event: make(chan int, 1),
    }

    register.Observer(Name, new)

    return new
}

func (r *Output) Init() {
    // 等待配置更新完成的信号
    <-r.event
    fmt.Println("output init")

    if v := outputs.Value; v != nil {
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

func (r *Output) Main() {
    return
}

func (r *Output) Exit(code int) {
    return
}

func (r *Output) Update(configure prototype.Object) int {
    if configure == nil {
        return state.Error
    }

    exist := true
    outputs.Value, exist = configure.(map[interface{}]interface{})[Name]
    if !exist {
        fmt.Println("Not found inputs configure")
        return state.Error
    }

    r.event <- 1
    return state.Ok
}

func (r *Output) Load(m module.Template) {
    r.children = append(r.children, m)
}

func init() {
    register.Module(module.Worker, Name, commands, New)
}
