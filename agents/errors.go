package agents

import "errors"

var (
    ErrType      = errors.New("Type is not equal")
    ErrConfigure = errors.New("Not found agents configure")
)
