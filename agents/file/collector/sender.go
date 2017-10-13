package collector

import (
	"github.com/rookie-xy/hubble/event"
)

func (c *Collector) Publish(event event.Event) bool {
	/*
    if err := c.client.Sender(event, false); err != nil {
    	fmt.Errorf("send client error", err)
        return false
    }
	*/
/*
    if err := c.sincedb.Sender(data.GetEvent(), false); err != nil {
    	fmt.Errorf("send sincedb error", err)
        return false
    }
*/
    return true
}
