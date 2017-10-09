package data

import (
    "testing"

    "github.com/rookie-xy/modules/agents/file/state"
    "github.com/stretchr/testify/assert"
)

func TestNewData(t *testing.T) {
	data := New()

//	assert.False(t, data.HasEvent())
	assert.False(t, data.HasState())

	data.Set(state.State{Source: "-"})

//	assert.False(t, data.HasEvent())
	assert.True(t, data.HasState())

//	data.Event.Fields = common.MapStr{}

//	assert.True(t, data.HasEvent())
	assert.True(t, data.HasState())
}

func TestGetEvent(t *testing.T) {
	//data := NewData()
//	data.Event.Fields = common.MapStr{"hello": "world"}
//	out := common.MapStr{"hello": "world"}
//	assert.Equal(t, out, data.GetEvent().Fields)
}
