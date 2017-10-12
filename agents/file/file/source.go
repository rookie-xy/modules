package file

import (
	"os"
	"github.com/rookie-xy/hubble/source"
	"github.com/modules/agents/log/open"
	"fmt"
	"github.com/rookie-xy/modules/agents/file/state"
	"errors"
)

type Source struct {
	*os.File
    state state.State
}

func (Source) Continuable() bool { return true }

func ReadOpen(path string) (*os.File, error) {
    flag := os.O_RDONLY
    perm := os.FileMode(0)
    return os.OpenFile(path, flag, perm)
}

func New(state state.State) (*Source, error) {
	f, err := ReadOpen(state.Source)
	if err != nil {
		return nil, fmt.Errorf("Failed opening %s: %s", state.Source, err)
	}

    source := &Source{
    	File: f,
    	state: state,
    }

    return source, nil
}

func (s *Source) Init() error {
	info, err := s.Stat()
	if err != nil {
		return fmt.Errorf("Failed getting stats for file %s: %s", s.state.Source, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("Tried to file non regular file: %q %s", info.Mode(), info.Name())
	}

	// Compares the stat of the opened file to the state given by the prospector. Abort if not match.
	if !os.SameFile(s.state.Fileinfo, info) {
		return errors.New("file info is not identical with opened file. Aborting harvesting and retrying file later again")
	}

	// get file offset. Only update offset if no error
	offset, err := s.initFileOffset(s.File)
	if err != nil {
		return err
	}

	fmt.Printf("collector Setting offset for file: %s. Offset: %d ", s.state.Source, offset)
	s.state.Offset = offset

	return nil
}

func (s *Source) initFileOffset(file *os.File) (int64, error) {
	// continue from last known offset
	if s.state.Offset > 0 {
		fmt.Printf("collector Set previous offset for file: %s. Offset: %d ", s.state.Source, s.state.Offset)
		return file.Seek(s.state.Offset, os.SEEK_SET)
	}

	// get offset from file in case of encoding factory was required to read some data.
	fmt.Printf("collector Setting offset for file based on seek: %s", s.state.Source)
	return file.Seek(0, os.SEEK_CUR)
}
