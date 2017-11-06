package collector

import (
	"github.com/rookie-xy/hubble/event"
	"fmt"
)

func (c *Collector) Publish(event event.Event) bool {
    if err := c.output.Sender(event); err != nil {
        return false
	}

	if c.conf.Client {
		fmt.Println("88888888888888888888888888888888888888")
//		return c.sincedb.Sender(event) == nil
	}

    return true
}
