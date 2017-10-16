package event

import (
    "github.com/rookie-xy/modules/agents/file/scanner"
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/event"
    "github.com/rookie-xy/modules/agents/file/state"
)

type Event struct {
    File      state.State
    Header    types.Map        `json:"header"`
    Body     *scanner.Message  `json:"body"`
    Footer    Footer           `json:"footer"`
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

func (e *Event) GetHeader() types.Map {
    return e.Header
}

func (e *Event) GetBody() event.Message {
	return e.Body
}

func (e *Event) GetFooter() []byte {
    return nil
}
