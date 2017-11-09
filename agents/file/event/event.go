package event

import (
    "github.com/rookie-xy/modules/agents/file/scanner"
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/models/file"
)

type Event struct {
    State    file.State
    Header   types.Map       `json:"header"`
    Body   *scanner.Message  `json:"body"`
    Footer   Footer          `json:"footer"`
}

type Footer struct {
    CheckSum  []byte  `json:"checksum"`
}

func New() *Event {
    return &Event{}
}

func (e *Event) ID() string {
    return ""
}

func (e *Event) GetState() file.State {
    return e.State
}

func (e *Event) GetHeader() types.Map {
    return e.Header
}

func (e *Event) GetBody() adapter.MessageEvent {
	return e.Body
}

func (e *Event) GetFooter() []byte {
    return nil
}


func (e *Event) Off() {

}

func (e *Event) On() bool {
    return true
}
