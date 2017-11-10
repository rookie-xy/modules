package collector

import (
	"github.com/rookie-xy/hubble/event"
)

func (c *Collector) Publish(event event.Event) bool {
	if c.client {
	    if err := c.output.Sender(event); err != nil {
            return false
	    }

        return c.sinceDB.Sender(event) == nil

	}

    if err := c.output.Sender(event); err != nil {
        return false
    }

    return true
}
