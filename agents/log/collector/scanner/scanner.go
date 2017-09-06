package text

import (
    "sync"

    "github.com/satori/go.uuid"

    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/types"

    "github.com/rookie-xy/modules/agents/log/util"
    "github.com/rookie-xy/modules/agents/log/file/state"
    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/hubble/proxy"
)

type Scanner struct {
    id       uuid.UUID
    //config   config
    source   Source // the source being watched

    // shutdown handling
    done      chan struct{}
    stopOnce  sync.Once
    stopWg   *sync.WaitGroup

    // internal harvester state
    state     state.State
    states   *state.States
    log       log.Log
    codec     codec.Codec
    client    proxy.Forward
/*
    // file reader pipeline
    reader          reader.Reader
    encodingFactory encoding.EncodingFactory
    encoding        encoding.Encoding

    // event/state publishing
    forwarder    *harvester.Forwarder
    publishState func(*util.Data) bool
    */
}

type Scanner struct {

}

func New(log log.Log) *Text {
    return &Text{
        log: log,
    }
}

func (c *Text) Init(group, Type string,
                              codec, client types.Value) error {
    return nil
}

func (c *Text) Setup(codec types.Value) error {
    return nil
}

func (c *Text) Scanner(fstates *state.States) *Scanner {
    return nil
}

// SendStateUpdate send an empty event with the current state to update the registry
// close_timeout does not apply here to make sure a collector is closed properly. In
// case the output is blocked the collector will stay open to make sure no new collector
// is started. As soon as the output becomes available again, the finished state is written
// and processing can continue.
func (r *Text) Update(fs state.State) {
    if !r.source.HasState() {
        return
    }

    fmt.Println("collector", "Update state: %s, offset: %v", r.state.Source, r.state.Offset)
    r.states.Update(r.state)

    d := util.NewData()
    d.SetState(r.state)
    //h.publishState(d)
}

func (c *Text) ID() uuid.UUID {
    return 1
}

func (c *Text) Run() error {
    return nil
}

func (c *Text) Stop() {
    return
}
