package source

import (
    "fmt"
    "os"
    "errors"

  . "github.com/rookie-xy/hubble/log/level"

    "github.com/rookie-xy/hubble/log"
    "github.com/rookie-xy/hubble/models/file"
)

type Source struct {
	*os.File

    state  file.State
    log    log.Factory
}

func (Source) Continuable() bool { return true }

func New(state file.State, log log.Factory) (*Source, error) {
	f, err := Open(state.Source)
	if err != nil {
		return nil, fmt.Errorf("Failed opening %s: %s", state.Source, err)
	}

    source := &Source{
    	state: state,
    	log:   log,
    }

    if err := source.validate(f); err != nil {
    	return nil, err
	}
	source.File = f

    return source, nil
}

func (s *Source) validate(f *os.File) error {
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("Failed getting stats for source %s: %s", s.state.Source, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("Tried to source non regular source: %q %s", info.Mode(), info.Name())
	}

	// Compares the stat of the opened source to the models given by the prospector. Abort if not match.
	if !os.SameFile(s.state.Fileinfo, info) {
		return errors.New("Source info is not identical with opened source. Aborting harvesting and retrying source later again")
	}

	// get source offset. Only update offset if no error
	offset, err := s.offset(f)
	if err != nil {
		return err
	}

	s.log(DEBUG, "Collector setting offset for source: %s. offset: %d", s.state.Source, offset)
	s.state.Offset = offset

	return nil
}

func (s *Source) offset(file *os.File) (int64, error) {
	// continue from last known offset
	if s.state.Offset > 0 {
        s.log(DEBUG,"Collector set previous offset for source: %s. Offset: %d", s.state.Source, s.state.Offset)
		return file.Seek(s.state.Offset, os.SEEK_SET)
	}

	// get offset from source in case of encoding factory was required to read some data.
	s.log(DEBUG,"Collector setting offset for source based on seek: %s\n", s.state.Source)
	return file.Seek(0, os.SEEK_CUR)
}
