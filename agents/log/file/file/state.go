package file

import (
    "os"
    "time"
)

type State struct {
    Id         string        `json:"-"` // local unique id to make comparison more efficient
    Finished   bool          `json:"-"` // harvester state
    Fileinfo   os.FileInfo   `json:"-"` // the file info
    Source     string        `json:"source"`
    Offset     int64         `json:"offset"`
    Timestamp  time.Time     `json:"timestamp"`
    TTL        time.Duration `json:"ttl"`
    Type       string        `json:"type"`
    Key        Key
}

// NewState creates a new file state
func NewState(fi os.FileInfo, path string, from string) State {
    return State{
        Finished:    false,
        Fileinfo:    fi,
        Source:      path,
        Timestamp:   time.Now(),
        TTL:         -1, // By default, state does have an infinite ttl
        Type:        from,
        Key: GetOSState(fi),
    }
}

// ID returns a unique id for the state as a string
func (r *State) ID() string {
    // Generate id on first request. This is needed as id is
    // not set when converting back from json
    if r.Id == "" {
        r.Id = r.Key.String()
				}

    return r.Id
}

// IsEqual compares the state to an other state supporing
// stringer based on the unique string
func (r *State) IsEqual(c *State) bool {
    return r.ID() == c.ID()
}

// IsEmpty returns true if the state is empty
func (r *State) IsEmpty() bool {
    return *r == State{}
}
