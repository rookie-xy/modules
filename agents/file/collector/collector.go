package collector

import (
    "io"
   	"fmt"
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
    log       log.Log

    once      sync.Once

    done     chan struct{}
    client    bool
}

func New(log log.Log) *Collector {
    return &Collector{
        log: log,
        id:  uuid.NewV4(),
        client: true,
        done: make(chan struct{}),
    }
}

func (c *Collector) Init(input input.Input, codec codec.Codec, state file.State,
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
    if err := scanner.Init(codec, state); err != nil {
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
				//c.log.Info("End of source reached: %s. Closing because close_eof is enabled.", c.file.Source)
				fmt.Printf("End of source reached: %s. Closing because close_eof is enabled.\n", c.state.Source)
			case ErrClosed:
				//c.log.Info("Reader was closed: %s. Closing.", c.file.Source)
				fmt.Printf("Reader was closed: %s. Closing.\n", c.state.Source)
			case ErrRemoved:
                //c.log.Info("File was removed: %s. Closing because close_removed is enabled.", c.file.Source)
                fmt.Printf("File was removed: %s. Closing because close_removed is enabled.\n", c.state.Source)
			case ErrRenamed:
				//c.log.Info("File was renamed: %s. Closing because close_renamed is enabled.", c.file.Source)
				fmt.Printf("File was renamed: %s. Closing because close_renamed is enabled.\n", c.state.Source)
			case ErrTooLong:
				fmt.Printf("File was too long: %s.\n", c.state.Source)
            case ErrInactive:
            	//c.log.Info("File is inactive: %s. Closing because close_inactive of %v reached.", c.file.Source, c.config.CloseInactive)
            	fmt.Printf("File is inactive: %s. Closing because close_inactive of %v reached.\n", c.state.Source, c.conf.Expire)
			case ErrFinalToken:
				fmt.Printf("File was FinalToken: %s.\n", c.state.Source)
			case ErrFileTruncate:
                //c.log.Info("File was truncated. Begin reading source from offset 0: %s", c.file.Source)
				c.state.Offset = 0
				fmt.Printf("File was truncated. Begin reading source from offset 0: %s\n", c.state.Source)
			case ErrAdvanceTooFar:
				fmt.Printf("File was AdvanceTooFar: %s\n", c.state.Source)
			case ErrNegativeAdvance:
				fmt.Printf("File was NegativeAdvance: %s\n", c.state.Source)
			default:
                //c.log.Err("Read line error: %s; File: ", c.scanner.Err(), c.file.Source)
                fmt.Printf("Read line error: %s; File: %s\n", c.scanner.Err(), c.state.Source)
            }

            return nil
	    }

        if c.state.Offset == 0 {
            message.Content = bytes.Trim(message.Content, "\xef\xbb\xbf")
        }

        state := c.getState()
        state.Offset += int64(message.Bytes) + 1 // add one because \n
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

        if !c.Publish(event) {
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

func (c *Collector) getState() file.State {
    state := c.state
	// refreshes the values in State with the values from the collector itself
	state.ID = file.Id(c.state.Fileinfo)
	return state
}

func (c *Collector) Update(fs file.State) {
    fmt.Printf("collector update state: %s, offset: %v\n", c.state.Source, c.state.Offset)
    c.states.Update(fs)
}

func (c *Collector) clean() {
	// Mark collector as finished
	c.state.Finished = true

	fmt.Printf("collector stopping collector for file: %s\n", c.state.Source)
	defer fmt.Printf("collector collector cleanup finished for file: %s\n", c.state.Source)

	// Make sure file is closed as soon as collector exits
	// If file was never opened, it can't be closed
	if c.source != nil {

		// close file handler
		c.source.Close()

		fmt.Printf("collector Closing file: %s\n", c.state.Source)

		// On completion, push offset so we can continue where we left off if we relaunch on the same file
		// Only send offset if file object was created successfully
		c.Update(c.state)
	} else {
		fmt.Printf("Stopping collector, NOT closing file as file info not available: %s\n", c.state.Source)
	}
}
