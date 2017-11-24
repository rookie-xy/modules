package collector

import (
	"github.com/rookie-xy/modules/agents/file/event"
	"fmt"
)

func (c *Collector) Publish(event *event.Event) bool {
	c.states.Update(event.Footer)

	if c.client {
	    if err := c.output.Sender(event); err != nil {
            return false
	    }

        return c.sinceDB.Sender(event) == nil

	}

    if err := c.output.Sender(event); err != nil {
        return false
    } else {
    	fmt.Println("SENDERRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRROKKKKKKKKKKKKK")
	}

    return true
}
