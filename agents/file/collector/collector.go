package collector

import (
    "io"
	"os"
   	"fmt"
    "errors"
    "bytes"

    "github.com/satori/go.uuid"

    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/hubble/proxy"
    "github.com/rookie-xy/hubble/command"
    "github.com/rookie-xy/hubble/factory"
    "github.com/rookie-xy/hubble/valve"

    "github.com/rookie-xy/modules/agents/file/state"
	"github.com/rookie-xy/modules/agents/file/source"
	"github.com/rookie-xy/modules/agents/file/scanner"
	"github.com/rookie-xy/modules/agents/file/event"
	"github.com/rookie-xy/modules/agents/file/data"
)

type Collector struct {
    id       uuid.UUID
    source   source.Source

    // internal state
    state      state.State
    states    *state.States
    log        log.Log
    codec      codec.Codec
    client     proxy.Forward
    sincedb    proxy.Forward
    scanner   *scanner.Scanner
    valve      valve.Valve
}

func New(log log.Log) *Collector {
    return &Collector{
        log: log,
        id:  uuid.NewV4(),
    }
}

func (c *Collector) Init(group, Type string,
                         codec, client *command.Command) error {
                         	/*
    event := message.New()
    if err := event.SetHeader("group", group); err != nil {
        return err
    }
    if err := event.SetHeader("type", Type); err != nil {
        return err
    }
    c.message = event
                         	*/

	key := codec.GetFlag() + "." + codec.GetKey()
    if codec, err := factory.Codec(key, c.log, codec.GetValue()); err != nil {
        return err
    } else {
        c.codec = codec
    }

 	key = client.GetFlag() + "." + client.GetKey()
    if client, err := factory.Client(key, c.log, client.GetValue()); err != nil {
        return err
    } else {
        c.client = client
    }

    if client, err := factory.Forward("plugin.client.sincedb"); err != nil {
        return err
    } else {
        c.sincedb = client
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
				c.log.Info("End of file reached: %s. Closing because close_eof is enabled.", c.state.Source)
			case source.ErrClosed:
				c.log.Info("Reader was closed: %s. Closing.", c.state.Source)
			case source.ErrRemoved:
                c.log.Info("File was removed: %s. Closing because close_removed is enabled.", c.state.Source)
			case source.ErrRenamed:
				c.log.Info("File was renamed: %s. Closing because close_renamed is enabled.", c.state.Source)
			case source.ErrTooLong:
            case source.ErrInactive:
            	c.log.Info("File is inactive: %s. Closing because close_inactive of %v reached.", c.state.Source, c.config.CloseInactive)
			case source.ErrFinalToken:
			case source.ErrFileTruncate:
                c.log.Info("File was truncated. Begin reading file from offset 0: %s", c.state.Source)
				c.state.Offset = 0
			case source.ErrAdvanceTooFar:
			case source.ErrNegativeAdvance:
			default:
                c.log.Err("Read line error: %s; File: ", c.scanner.Err(), c.state.Source)
            }

            return nil
	    }

        if c.state.Offset == 0 {
            message.Content = bytes.Trim(message.Content, "\xef\xbb\xbf")
        }

        state := c.getState()
        state.Offset += int64(message.Bytes)

        data := data.New()
   		if c.source.HasState() {
			data.Set(state)
		}

        if !message.IsEmpty() /*&& c.valve.Filter(c.scanner.Text())*/ {
        	header := event.Header{
        		"group": c.group,
        		"type":  c.Type,
			}

			data.Event = event.Event{
				Magic: 0x0100,
				Header: header,
				Body: message,
			}
            /*
			if c.fingerprint {
				md5 := md5.New()
				event.Footer{
					CheckSum: md5.Sum([]byte(data.Event))
				}
			}
            */

		}

		if !c.Sender(data) {
		    return nil
		}

		c.state = state
    }
}

func (c *Collector) Sender(data *data.Data) bool {
    if err := c.client.Sender(data.Event, false); err != nil {
    	fmt.Errorf("send client error", err)
        return false
    }

    if err := c.sincedb.Sender(data.GetEvent(), false); err != nil {
    	fmt.Errorf("send sincedb error", err)
        return false
    }

    return true
}

func (c *Collector) Stop() {

}

func (c *Collector) Setup() error {
    if err := c.openFile(); err != nil {
        return err
    }

    log := source.New(c.log)
	if err := log.Init(nil); err != nil {
		return err
	}

    scanner := scanner.New(log)
    if err := scanner.Init(c.codec); err != nil {
    	return err
	}
    c.scanner = scanner

    return nil
}

func (c *Collector) getState() state.State {
	if !c.source.HasState() {
		return state.State{}
	}

	state := c.state

	// refreshes the values in State with the values from the harvester itself
	//state.FileStateOS = file.GetOSState(h.state.Fileinfo)
	return state
}

func ReadOpen(path string) (*os.File, error) {
    flag := os.O_RDONLY
    perm := os.FileMode(0)
    return os.OpenFile(path, flag, perm)
}

type File struct {
	*os.File
}

func (File) Continuable() bool { return true }
func (File) HasState() bool    { return true }

func (c *Collector) openFile() error {
	f, err := ReadOpen(c.state.Source)
	if err != nil {
		return fmt.Errorf("Failed opening %s: %s", c.state.Source, err)
	}

	// Makes sure file handler is also closed on errors
	err = c.validateFile(f)
	if err != nil {
		f.Close()
		return err
	}

	c.source = File{File: f}
	return nil
}

func (c *Collector) validateFile(f *os.File) error {
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("Failed getting stats for file %s: %s", c.state.Source, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("Tried to open non regular file: %q %s", info.Mode(), info.Name())
	}

	// Compares the stat of the opened file to the state given by the prospector. Abort if not match.
	if !os.SameFile(c.state.Fileinfo, info) {
		return errors.New("file info is not identical with opened file. Aborting harvesting and retrying file later again")
	}
/*
	c.encoding, err = c.encodingFactory(f)
	if err != nil {

		if err == transform.ErrShortSrc {
			fmt.Printf("Initialising encoding for '%v' failed due to file being too short", f)
		} else {
			fmt.Printf("Initialising encoding for '%v' failed: %v", f, err)
		}
		return err
	}
*/

	// get file offset. Only update offset if no error
	offset, err := c.initFileOffset(f)
	if err != nil {
		return err
	}

	fmt.Printf("collector Setting offset for file: %s. Offset: %d ", c.state.Source, offset)
	c.state.Offset = offset

	return nil
}

func (c *Collector) initFileOffset(file *os.File) (int64, error) {
	// continue from last known offset
	if c.state.Offset > 0 {
		fmt.Printf("collector Set previous offset for file: %s. Offset: %d ", c.state.Source, c.state.Offset)
		return file.Seek(c.state.Offset, os.SEEK_SET)
	}

	// get offset from file in case of encoding factory was required to read some data.
	fmt.Printf("collector Setting offset for file based on seek: %s", c.state.Source)
	return file.Seek(0, os.SEEK_CUR)
}

func (r *Collector) Update(fs state.State) {
    if !r.source.HasState() {
        return
    }

    fmt.Println("collector update state: %s, offset: %v", r.state.Source, r.state.Offset)
    r.states.Update(r.state)

//    d := data.NewData()
//    d.SetState(r.state)
    //h.publishState(d)
}
