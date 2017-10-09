package scanner

import (
    "bufio"
    "github.com/rookie-xy/modules/agents/file/source"
    "github.com/rookie-xy/hubble/codec"
    "github.com/rookie-xy/modules/agents/file/state"
)

type Scanner struct {
    *bufio.Scanner
    id uint64
}

func New(log *source.Log) *Scanner {
    return nil
}

func (s *Scanner) Init(codec codec.Codec, state state.State) error {
    s.id = state.Lno
    s.Split(codec.Decode)
    return nil
}

func (s *Scanner) ID() uint64 {
    s.id++
    return s.id
}

func (s *Scanner) Scan() (*Message, bool) {
    if s.Scanner.Scan() {
        message := &Message{
            Id:      s.ID(),
            Content: s.Bytes(),
        }

        message.Bytes = len(message.Content)
        return message, true
    }

    return nil, false
}
