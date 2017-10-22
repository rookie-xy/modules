package collector

import (
	"github.com/rookie-xy/hubble/event"
)

func (c *Collector) Publish(event event.Event) bool {
    if err := c.output.Sender(event); err != nil {
        return false
	}

	if c.conf.Client {
		return c.sincedb.Sender(event) == nil
	}

    return true
}
