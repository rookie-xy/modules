package collector

import (
    "io"
    "bytes"

    "github.com/satori/go.uuid"

    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/filter"
	"github.com/rookie-xy/hubble/types"
	"github.com/rookie-xy/hubble/input"
  . "github.com/rookie-xy/hubble/source"
	"github.com/rookie-xy/hubble/output"
	"github.com/rookie-xy/hubble/models/file"

	"github.com/rookie-xy/modules/agents/file/scanner"
	"github.com/rookie-xy/modules/agents/file/event"
    "github.com/rookie-xy/modules/agents/file/configure"
	"github.com/rookie-xy/hubble/proxy"
	"github.com/rookie-xy/hubble/factory"
	"github.com/rookie-xy/modules/agents/file/source"
	"github.com/rookie-xy/hubble/types/value"
	"github.com/rookie-xy/hubble/codec"
	"sync"
  . "github.com/rookie-xy/hubble/log/level"
	"github.com/rookie-xy/hubble/adapter"
	"github.com/rookie-xy/modules/agents/file/utils"
)

type Collector struct {
    id        uuid.UUID

    state     file.State
    states   *file.States

    conf     *configure.Configure

    source    Source
    input     input.Input
    scanner  *scanner.Scanner
    filter    filter.Filter
    output    proxy.Forward

    sinceDB   proxy.Forward
    log.Log
    level     Level

    once      sync.Once

    done      chan struct{}
    client    bool
}

func New(log log.Log) *Collector {
    return &Collector{
        Log: log,
        level: adapter.ToLevelLog(log).Get(),
        id:  uuid.NewV4(),
        client: true,
        done: make(chan struct{}),
    }
}

func (c *Collector) Init(input input.Input, decoder codec.Decoder, state file.State,
	                     states *file.States, conf *configure.Configure) error {
	var err error
    source, err := source.New(state, c.log)
    if err != nil {
        return err
    }

    if err := input.Init(source); err != nil {
    	return err
	}

    scanner := scanner.New(input)
    if err := scanner.Init(decoder, state); err != nil {
    	return err
	}

    c.conf    = conf
	c.state   = state
	c.states  = states
	c.source  = source
	c.input   = input
    c.scanner = scanner

    if c.conf.Output != nil {
		pluginName := c.conf.Output.GetFlag() + "." + c.conf.Output.GetKey()
		c.output, err = factory.Output(pluginName, c.Log, c.conf.Output.GetValue())
		if err != nil {
			return err
		}

		c.client = false

	} else {
	    key := c.conf.Client.GetFlag() + "." + c.conf.Client.GetKey()
        c.output, err = factory.Client(key, c.Log, c.conf.Client.GetValue())
        if err != nil {
            return err
        }

        key = c.conf.SinceDB.GetFlag() + "." + output.Name + "." + "sinceDB"
        c.sinceDB, err = factory.Output(key, c.Log, value.New(c.conf.SinceDB.GetKey()))
        if err != nil {
            return err
        }
	}

    return nil
}

func (c *Collector) ID() uuid.UUID {
    return c.id
}

func (c *Collector) Run() error {
	defer func() {
		c.Stop()
		c.clean()
	}()

    for {
        select {
        case <-c.done:
            return nil
        default:
		}

        message, keep := c.scanner.Scan()
        if !keep {
            switch c.scanner.Err() {
			case io.EOF:
				c.log(ERROR,"End of source reached: %s. Closing because close_eof is enabled",
					             c.state.Source)
			case ErrClosed:
				c.log(ERROR,"Reader was closed: %s. Closing.\n", c.state.Source)
			case ErrRemoved:
                c.log(ERROR,"File was removed: %s. Closing because close_removed is enabled",
                	             c.state.Source)
			case ErrRenamed:
				c.log(ERROR,"File was renamed: %s. Closing because close_renamed is enabled",
					             c.state.Source)
			case ErrTooLong:
				c.log(ERROR,"File was too long: %s", c.state.Source)
            case ErrInactive:
            	c.log(ERROR,"File is inactive: %s. Closing because close_inactive of %v reached",
            		             c.state.Source, c.conf.Expire)
			case ErrFinalToken:
				c.log(ERROR,"File was FinalToken: %s", c.state.Source)
			case ErrFileTruncate:
				c.state.Offset = 0
				c.log(ERROR,"File was truncated. Begin reading source from offset 0: %s",
					             c.state.Source)
			case ErrAdvanceTooFar:
				c.log(ERROR,"File was AdvanceTooFar: %s", c.state.Source)
			case ErrNegativeAdvance:
				c.log(ERROR,"File was NegativeAdvance: %s", c.state.Source)
			default:
                c.log(ERROR,"Read line error: %s; File: %s", c.scanner.Err(), c.state.Source)
            }
            return nil
	    }

        if c.state.Offset == 0 {
            message.Content = bytes.Trim(message.Content, "\xef\xbb\xbf")
        }

        state := utils.GetState(c.state)
        state.Offset += int64(message.Bytes) + 1 // add one because "\n"
        state.Lno = message.ID()

        event := &event.Event{
            Footer: state,
		}

        if !message.IsEmpty() /*&& c.filter.Handler(string(message.Bytes))*/ {
			event.Header = types.Map{
        		"group": c.conf.Group,
        		"type":  c.conf.Type,
			}
			event.Body = message
		}

        if !c.publish(event) {
		    return nil
		}
		c.state = state
    }
}

func (c *Collector) Stop() {
    c.output.Close()
    c.once.Do(func() {
        close(c.done)
    })
}

func (c *Collector) Update(fs file.State) {
    c.log(DEBUG,"collector update state: %s, offset: %v", c.state.Source, c.state.Offset)
    c.states.Update(fs)
}

func (c *Collector) publish(e *event.Event) bool {
	c.states.Update(e.Footer)

	if err := c.output.Sender(e); err != nil {
    	c.log(ERROR, "send event failure: %s", err)
        return false
    }

    c.log(DEBUG, "send event successful")

    if !c.client {
    	return true
	}
    return c.sinceDB.Sender(e) == nil
}

func (c *Collector) clean() {
	c.state.Finished = true

	c.log(DEBUG,"collector stopping collector for file: %s", c.state.Source)
	defer c.log(DEBUG,"collector collector cleanup finished for file: %s\n", c.state.Source)

	if c.source != nil {
		c.source.Close()
        c.log(DEBUG,"collector Closing file: %s\n", c.state.Source)
		c.Update(c.state)
	} else {
		c.log(DEBUG,"Stopping collector, NOT closing file as file info not available: %s\n",
			             c.state.Source)
	}
}

func (c *Collector) log(l Level, fmt string, args ...interface{}) {
    log.Print(c.Log, c.level, l, fmt, args...)
}
