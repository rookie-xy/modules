package collector

import (
	"os"
   	"fmt"
    "errors"
    "bytes"
    "bufio"

    "github.com/satori/go.uuid"

    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/hubble/proxy"

    "github.com/rookie-xy/modules/agents/file/state"
	"github.com/rookie-xy/hubble/command"
	"github.com/rookie-xy/hubble/factory"
	"github.com/rookie-xy/modules/agents/file/message"
	"github.com/rookie-xy/hubble/event"
	"github.com/rookie-xy/modules/agents/file/source"
)

type Collector struct {
    id       uuid.UUID
    source   source.Source

    // internal state
    state      state.State
    states    *state.States
    log        log.Log
    Log       *source.Log
    codec      codec.Codec
    client     proxy.Forward
    sincedb    proxy.Forward
    message   *message.Message
}

func New(log log.Log) *Collector {
    return &Collector{
        log: log,
        id:  uuid.NewV4(),
    }
}

func (c *Collector) Init(group, Type string,
                         codec, client *command.Command) error {
    event := message.New()
    if err := event.SetHeader("group", group); err != nil {
        return err
    }
    if err := event.SetHeader("type", Type); err != nil {
        return err
    }
    c.message = event

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
    scanner := bufio.NewScanner(c.log)
	scanner.Split(c.codec.Decode)
	id := c.state.Lno

    for {
        keep := scanner.Scan()
        if !keep {
            switch scanner.Err() {
			case bufio.ErrTooLong:
			case bufio.ErrFinalToken:
			case bufio.ErrAdvanceTooFar:
			case bufio.ErrNegativeAdvance:
            }
	    }
	    id++

	    message := message.New()
	    message.Id = id
	    content := scanner.Bytes()
	    //message.Bytes = len(message.Content)
        if c.state.Offset == 0 {
            content = bytes.Trim(content, "\xef\xbb\xbf")
        }
        message.SetBody(content)

        state := c.getState()
        state.Offset += message.GetBodyLength()

        event := event.New()
   		if c.source.HasState() {
			event.Set(state)
		}

        text := scanner.Text()

        if !message.IsEmpty() {

		}

		if !c.Send(event) {
		    return nil
		}

		c.state = state
    }
}

func (c *Collector) Stop() {

}

func (c *Collector) Send(event event.Event) bool {
    if err := c.client.Sender(event, false); err != nil {
    	fmt.Errorf("send client error", err)
        return false
    }

    if err := c.sincedb.Sender(event, false); err != nil {
    	fmt.Errorf("send sincedb error", err)
        return false
    }

    return true
}

func (c *Collector) Setup() error {
    if err := c.openFile(); err != nil {
        return err
    }

    log := source.New(c.log)
	if err := log.Init(nil); err != nil {
		return err
	}
	c.Log = log

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

//    d := util.NewData()
//    d.SetState(r.state)
    //h.publishState(d)
}
