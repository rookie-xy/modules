package collector

import (
    "io"
   	"fmt"
    "bytes"

    "github.com/satori/go.uuid"

    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/valve"

    "github.com/rookie-xy/modules/agents/file/state"
	"github.com/rookie-xy/modules/agents/file/scanner"
	"github.com/rookie-xy/modules/agents/file/event"
	"github.com/rookie-xy/hubble/types"
	"github.com/rookie-xy/modules/agents/file/configure"
	"github.com/rookie-xy/modules/agents/file/file"
	"github.com/rookie-xy/hubble/input"
	"github.com/rookie-xy/hubble/source"
)

type Collector struct {
   *configure.Configure

    id         uuid.UUID
    log        log.Log

    // internal state
    state  state.State
    states *state.States

    source  source.Source
    input   input.Input
    scanner *scanner.Scanner

    valve   valve.Valve
}

func New(log log.Log) *Collector {
    return &Collector{
        log: log,
        id:  uuid.NewV4(),
        //fingerprint: false,
    }
}

func (c *Collector) Init(input input.Input, state state.State) error {
	var err error
    source, err := file.New(state)
    if err != nil {
        return err
    }

    if err := source.Init(); err != nil {
    	return err
	}

    if err := input.Init(source); err != nil {
    	return err
	}

    scanner := scanner.New(input)
    if err := scanner.Init(c.Codec, c.state); err != nil {
    	return err
	}

	c.state   = state
	c.source  = source
	c.input   = input
    c.scanner = scanner

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
				//c.log.Info("End of file reached: %s. Closing because close_eof is enabled.", c.state.Source)
			case source.ErrClosed:
				//c.log.Info("Reader was closed: %s. Closing.", c.state.Source)
			case source.ErrRemoved:
                //c.log.Info("File was removed: %s. Closing because close_removed is enabled.", c.state.Source)
			case source.ErrRenamed:
				//c.log.Info("File was renamed: %s. Closing because close_renamed is enabled.", c.state.Source)
			case source.ErrTooLong:
            case source.ErrInactive:
            	//c.log.Info("File is inactive: %s. Closing because close_inactive of %v reached.", c.state.Source, c.config.CloseInactive)
			case source.ErrFinalToken:
			case source.ErrFileTruncate:
                //c.log.Info("File was truncated. Begin reading file from offset 0: %s", c.state.Source)
				c.state.Offset = 0
			case source.ErrAdvanceTooFar:
			case source.ErrNegativeAdvance:
			default:
                //c.log.Err("Read line error: %s; File: ", c.scanner.Err(), c.state.Source)
            }

            return nil
	    }

        if c.state.Offset == 0 {
            message.Content = bytes.Trim(message.Content, "\xef\xbb\xbf")
        }

        state := c.getState()
        state.Offset += int64(message.Bytes)

        event := event.Event{
            File: state,
		}

        if !message.IsEmpty() /*&& c.valve.Filter(c.scanner.Text())*/ {
			event.Header = types.SiMap{
        		"group": c.Configure.Group,
        		"type":  c.Configure.Type,
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

func (c *Collector) getState() state.State {
	state := c.state

	// refreshes the values in State with the values from the harvester itself
	//state.FileStateOS = file.GetOSState(h.state.Fileinfo)
	return state
}

func (r *Collector) Update(fs state.State) {
    fmt.Println("collector update state: %s, offset: %v", r.state.Source, r.state.Offset)
    r.states.Update(r.state)

//    d := data.NewData()
//    d.SetState(r.state)
    //h.publishState(d)
}


	/*
    if client, err := factory.Forward("plugin.client.sincedb"); err != nil {
        return err
    } else {
        c.sincedb = client
    }
	*/
