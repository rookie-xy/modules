package collector

import (
	"os"
   	"fmt"
    "sync"
    "errors"

    "github.com/satori/go.uuid"

    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/types"
    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/hubble/proxy"

    "github.com/rookie-xy/modules/agents/file/state"
)

type Collector struct {
    id       uuid.UUID
    source   Source

    // shutdown handling
    done      chan struct{}
    stopOnce  sync.Once
    stopWg   *sync.WaitGroup

    // internal state
    state      state.State
    states    *state.States
    log        log.Log
    codec      codec.Codec
    client     proxy.Forward
    sincedb    proxy.Reverse
}

func New(log log.Log) *Collector {
    return &Collector{
        log: log,
        id:  uuid.NewV4(),
    }
}

func (c *Collector) Init(group, Type string,
                         codec, client types.Value) error {
    return nil
}

func (c *Collector) ID() uuid.UUID {
    return c.id
}

func (c *Collector) Run() error {
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
    return nil
}

func (c *Collector) Stop() {

}

func (c *Collector) Setup() error {
    return nil
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
