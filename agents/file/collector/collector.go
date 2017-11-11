package collector

import (
    "io"
//   	"fmt"
    "bytes"

    "github.com/satori/go.uuid"

    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/filter"
	"github.com/rookie-xy/hubble/types"
	"github.com/rookie-xy/hubble/input"
     src "github.com/rookie-xy/hubble/source"
	"github.com/rookie-xy/hubble/output"
	"github.com/rookie-xy/hubble/models/file"

	"github.com/rookie-xy/modules/agents/file/scanner"
	"github.com/rookie-xy/modules/agents/file/event"
    "github.com/rookie-xy/modules/agents/file/configure"
	"github.com/rookie-xy/hubble/proxy"
	"github.com/rookie-xy/hubble/factory"
	"github.com/rookie-xy/modules/agents/file/source"
	"fmt"
	"github.com/rookie-xy/hubble/types/value"
)

type Collector struct {
    id        uuid.UUID

    state     file.State
    states   *file.States

    conf     *configure.Configure

    source    src.Source
    input     input.Input
    scanner  *scanner.Scanner
    filter    filter.Filter
    output    proxy.Forward

    sinceDB   proxy.Forward
    log       log.Log

    client    bool
}

func New(log log.Log) *Collector {
    return &Collector{
        log: log,
        id:  uuid.NewV4(),
        client: true,
        //fingerprint: false,
    }
}

func (c *Collector) Init(input input.Input, state file.State,
	                     states *file.States, conf *configure.Configure) error {
	var err error
    source, err := source.New(state)
    if err != nil {
        return err
    }

    if err := input.Init(source); err != nil {
    	return err
	}

    scanner := scanner.New(input)
    if err := scanner.Init(conf.Codec, state); err != nil {
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
		fmt.Printf("\nxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx: %s\n", pluginName)
		c.output, err = factory.Output(pluginName, c.log, c.conf.Output.GetValue())
		if err != nil {
			return err
		}

		c.client = false

	} else {
	    key := c.conf.Client.GetFlag() + "." + c.conf.Client.GetKey()
        c.output, err = factory.Client(key, c.log, c.conf.Client.GetValue())
        if err != nil {
            return err
        }

        key = c.conf.SinceDB.GetFlag() + "." + output.Name + "." + "sinceDB"
        c.sinceDB, err = factory.Output(key, c.log, value.New(c.conf.SinceDB.GetKey()))
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
    for {
        message, keep := c.scanner.Scan()
        if !keep {
            switch c.scanner.Err() {
			case io.EOF:
				//c.log.Info("End of source reached: %s. Closing because close_eof is enabled.", c.models.Source)
			case src.ErrClosed:
				//c.log.Info("Reader was closed: %s. Closing.", c.models.Source)
			case src.ErrRemoved:
                //c.log.Info("File was removed: %s. Closing because close_removed is enabled.", c.models.Source)
			case src.ErrRenamed:
				//c.log.Info("File was renamed: %s. Closing because close_renamed is enabled.", c.models.Source)
			case src.ErrTooLong:
            case src.ErrInactive:
            	//c.log.Info("File is inactive: %s. Closing because close_inactive of %v reached.", c.models.Source, c.config.CloseInactive)
			case src.ErrFinalToken:
			case src.ErrFileTruncate:
                //c.log.Info("File was truncated. Begin reading source from offset 0: %s", c.models.Source)
				c.state.Offset = 0
			case src.ErrAdvanceTooFar:
			case src.ErrNegativeAdvance:
			default:
                //c.log.Err("Read line error: %s; File: ", c.scanner.Err(), c.models.Source)
            }

            return nil
	    }

        if c.state.Offset == 0 {
            message.Content = bytes.Trim(message.Content, "\xef\xbb\xbf")
        }

        state := c.getState()
        state.Offset += int64(message.Bytes)

        event := &event.Event{
            Footer: state,
		}

        if !message.IsEmpty() /*&& c.filter.Handler(c.scanner.Text())*/ {
			event.Header = types.Map{
        		"group": c.conf.Group,
        		"type":  c.conf.Type,
			}
			event.Body = message
            /*
			if c.fingerprint {
				event.Footer = event.Footer{
				    CheckSum: nil,
				}
			}
            */
		}

        if !c.Publish(event) {
		    return nil
		}

		c.state = state
    }
}

func (c *Collector) Stop() {

}

func (c *Collector) getState() file.State {
    state := c.state
	// refreshes the values in State with the values from the harvester itself
	return state
}

func (c *Collector) Update(fs file.State) {
    fmt.Printf("collector update state: %s, offset: %v\n", c.state.Source, c.state.Offset)
    c.states.Update(fs)
}
