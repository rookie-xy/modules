package state

import (
    "os"
    "time"
)

type State struct {
    Id         string        `json:"-"` // local unique id to make comparison more efficient
    Finished   bool          `json:"-"` // harvester state
    Fileinfo   os.FileInfo   `json:"-"` // the file info
    Source     string        `json:"file"`
    Lno        uint64        `json:"lno"`
    Offset     int64         `json:"offset"`
    Timestamp  time.Time     `json:"timestamp"`
    TTL        time.Duration `json:"ttl"`
    Type       string        `json:"type"`
}

// New creates a new file state
func New() State {
    return State{
        Finished:    false,
        Timestamp:   time.Now(),
        TTL:         -1, // By default, state does have an infinite ttl
        Type:        "file",
    }
}

func (s State) Init(fi os.FileInfo, path, Type string) error {
    return nil
}

// ID returns a unique id for the state as a string
func (s State) ID() string {
    // Generate id on first request. This is needed as id is
    // not set when converting back from json
    if s.Id == "" {
        s.Id = s.Type
    }

    return s.Id
}

// IsEqual compares the state to an other state supporing
// stringer based on the unique string
func (s State) IsEqual(c *State) bool {
    return s.ID() == s.ID()
}

// IsEmpty returns true if the state is empty
func (s State) IsEmpty() bool {
    return s == State{}
}
