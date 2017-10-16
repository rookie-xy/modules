package events

import (
	"github.com/rookie-xy/hubble/log"
	"github.com/rookie-xy/hubble/types"
	"github.com/rookie-xy/hubble/event"
)

type Batch struct {
	log log.Log
	timeout int
}

func New(l log.Log) *Batch {
	return &Batch{
		log: l,
	}
}

func (b *Batch) Init(value types.Value) error {
	return nil
}

func (b *Batch) Enable() bool {
	return true
}

func (b *Batch) ID() string {
	return ""
}

func (b *Batch) GetHeader() types.Map {
	return nil
}

func (b *Batch) GetBody() event.Message {
	return nil
}

func (b *Batch) GetFooter() []byte {
	return nil
}

func (b *Batch) Put(event event.Event) int {
	return 0
}

func (b *Batch) Batch() []event.Event {
	return nil
}

func (b *Batch) Lengent() int {
	return 0
}
