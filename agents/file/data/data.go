package data

import (
	"github.com/rookie-xy/modules/agents/file/event"
	"github.com/rookie-xy/modules/agents/file/state"
)

type Data struct {
	Event event.Event
    state state.State
}

func New() *Data {
    return &Data{}
}

// SetState sets the state
func (d *Data) Set(state state.State) {
    d.state = state
}

// GetState returns the current state
func (d *Data) Get() state.State {
    return d.state
}

// HasState returns true if the data object contains state data
func (d *Data) HasState() bool {
    return d.state != state.State{}
}

// GetEvent returns the events in the data object
// In case meta data contains module and fileset data, the events is enriched with it
func (d *Data) GetEvent() event.Event {
    return d.Event
}
/*
// GetMetadata creates a common.MapStr containing the metadata to
// be associated with the events.
func (d *Data) GetMetadata() common.MapStr {
    return d.Event.Meta
}

// HasEvent returns true if the data object contains events data
func (d *Data) HasEvent() bool {
    return d.Event.Fields != nil
}
*/
