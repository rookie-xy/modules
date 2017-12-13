package collector

import (
	"github.com/rookie-xy/modules/agents/file/event"
  . "github.com/rookie-xy/hubble/log/level"
)

func (c *Collector) Publish(event *event.Event) bool {
	c.states.Update(event.Footer)

	if c.client {
	    if err := c.output.Sender(event); err != nil {
	    	c.log(ERROR, "client send failure: %s", err)
            return false
	    }

        return c.sinceDB.Sender(event) == nil
	}

    if err := c.output.Sender(event); err != nil {
    	c.log(ERROR, "output send failure: %s", err)
        return false
    } else {
        c.log(DEBUG, "output send successful")
	}

    return true
}
