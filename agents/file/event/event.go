package event

import (
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/adapter"
    "github.com/rookie-xy/hubble/models/file"
    "github.com/rookie-xy/modules/agents/file/message"
)

type Event struct {
    Header   types.Map       `json:"header"`
    Body    *message.Message `json:"body"`
    Footer   file.State      `json:"footer"`
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

func (e *Event) GetBody() adapter.MessageEvent {
	return e.Body
}

func (e *Event) GetFooter() file.State {
    return e.Footer
}

func (e *Event) Off() {

}

func (e *Event) On() bool {
    return true
}
