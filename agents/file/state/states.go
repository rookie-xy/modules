package state

import (
    "fmt"
    "sync"
    "time"
)

type States struct {
    States []State `json:"states"`
    sync.RWMutex
}

func News() *States {
    return &States{
        States: []State{},
    }
}

func (r *States) Update(newState State) {
    r.Lock()
    defer r.Unlock()

    index, _ := r.findPrevious(newState)
    newState.Timestamp = time.Now()

    if index >= 0 {
        r.States[index] = newState
    } else {
        r.States = append(r.States, newState)
        fmt.Println("finder", "New state added for %s", newState.Source)
    }
}

func (r *States) FindPrevious(newState State) State {
    r.RLock()
    defer r.RUnlock()

    _, state := r.findPrevious(newState)
    return state
}

func (s *States) findPrevious(newState State) (int, State) {
    for index, oldState := range s.States {
        if oldState.IsEqual(&newState) {
            return index, oldState
        }
    }

    return -1, State{}
}

func (r *States) Cleanup() int {
    r.Lock()
    defer r.Unlock()

    statesBefore := len(r.States)
    currentTime := time.Now()
    states := r.States[:0]

    for _, state := range r.States {
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

    r.States = states

    return statesBefore - len(r.States)
}

// Count returns number of states
func (r *States) Count() int {
    r.RLock()
    defer r.RUnlock()

    return len(r.States)
}

// Returns a copy of the file states
func (r *States) GetStates() []State {
    r.RLock()
    defer r.RUnlock()

    newStates := make([]State, len(r.States))
    copy(newStates, r.States)

    return newStates
}

func (r *States) SetStates(states []State) {
    r.Lock()
    defer r.Unlock()

    r.States = states
}

func (r *States) Copy() *States {
//    states := NewStates()
//    states.states = r.GetStates()

//    return states
    return nil
}
