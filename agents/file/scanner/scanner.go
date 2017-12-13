package scanner

import (
    "bufio"
    "errors"

    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/hubble/input"
    "github.com/rookie-xy/hubble/models/file"
    "github.com/rookie-xy/modules/agents/file/message"
    "github.com/rookie-xy/hubble/adapter"
)

type Scanner struct {
    *bufio.Scanner
    id      uint64
}

func New(s input.Input) *Scanner {
    scanner := &Scanner{
        Scanner: bufio.NewScanner(s),
        id: 0,
    }

    return scanner
}

func (s *Scanner) Init(decoder codec.Decoder, state file.State) error {
    if decoder == nil || state.IsEmpty() {
        return errors.New("state is empty or decoder is nil")
    }

    s.id = state.Lno
    s.Split(adapter.ToLogDecoder(decoder).LogDecode)
    return nil
}

func (s *Scanner) ID() uint64 {
    s.id++
    return s.id
}

func (s *Scanner) Scan() (*message.Message, bool) {
    if s.Scanner.Scan() {
        message := &message.Message{
            Id:      s.ID(),
            Content: s.Bytes(),
        }

        message.Bytes = len(message.Content)
        return message, true
    }

    return nil, false
}
