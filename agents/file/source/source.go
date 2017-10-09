package source

import (
    "io"
    "os"
    "errors"
    "bufio"
)

var (
    ErrFileTruncate    = errors.New("detected file being truncated")
    ErrRenamed         = errors.New("file was renamed")
    ErrRemoved         = errors.New("file was removed")
    ErrInactive        = errors.New("file inactive")
    ErrClosed          = errors.New("reader closed")

    ErrTooLong         = bufio.ErrTooLong
    ErrFinalToken      = bufio.ErrFinalToken
    ErrAdvanceTooFar   = bufio.ErrAdvanceTooFar
    ErrNegativeAdvance = bufio.ErrNegativeAdvance
)

type Source interface {
    io.ReadCloser

    Name() string
    Stat() (os.FileInfo, error)
    Continuable() bool // can we continue processing after EOF?
    HasState() bool    // does this source have a state?
}
