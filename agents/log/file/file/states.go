package file

import (
    "fmt"
    "sync"
    "time"
)

// States handles list of FileState
type States struct {
    states []State
    sync.RWMutex
}

func NewStates() *States {
    return &States{
        states: []State{},
    }
}

// Update updates a state. If previous state didn't exist, new one is created
func (r *States) Update(newState State) {
    r.Lock()
    defer r.Unlock()

    index, _ := r.findPrevious(newState)
    newState.Timestamp = time.Now()

    if index >= 0 {
        r.states[index] = newState
    } else {
		      // No existing state found, add new one
        r.states = append(r.states, newState)
				    fmt.Println("prospector", "New state added for %s", newState.Source)
    }
}

func (r *States) FindPrevious(newState State) State {
    r.RLock()
    defer r.RUnlock()

    _, state := r.findPrevious(newState)
    return state
}

// findPreviousState returns the previous state fo the file
// In case no previous state exists, index -1 is returned
func (s *States) findPrevious(newState State) (int, State) {
    // TODO: This could be made potentially more performance by using an index (harvester id) and only use iteration as fall back
    for index, oldState := range s.states {
        // This is using the FileStateOS for comparison as FileInfo
				    // identifiers can only be fetched for existing files
        if oldState.IsEqual(&newState) {
            return index, oldState
        }
    }

    return -1, State{}
}

// Cleanup cleans up the state array. All states which are older then `older` are removed
// The number of states that were cleaned up is returned
func (r *States) Cleanup() int {
    r.Lock()
    defer r.Unlock()

    statesBefore := len(r.states)
    currentTime := time.Now()
    states := r.states[:0]

    for _, state := range r.states {
        expired := (state.TTL > 0 && currentTime.Sub(state.Timestamp) > state.TTL)

        if state.TTL == 0 || expired {
            if state.Finished {
                fmt.Println("state", "State removed for %v because of older: %v", state.Source, state.TTL)
                continue // drop state
            } else {
												    fmt.Println("State for %s should have been dropped, but couldn't as state is not finished.", state.Source)
            }
        }

				    states = append(states, state) // in-place copy old state
    }

    r.states = states

    return statesBefore - len(r.states)
}

// Count returns number of states
func (r *States) Count() int {
    r.RLock()
    defer r.RUnlock()

    return len(r.states)
}

// Returns a copy of the file states
func (r *States) GetStates() []State {
    r.RLock()
    defer r.RUnlock()

    newStates := make([]State, len(r.states))
    copy(newStates, r.states)

    return newStates
}

// SetStates overwrites all internal states with the given states array
func (r *States) SetStates(states []State) {
    r.Lock()
    defer r.Unlock()

    r.states = states
}

// Copy create a new copy of the states object
func (r *States) Copy() *States {
    states := NewStates()
    states.states = r.GetStates()

    return states
}
