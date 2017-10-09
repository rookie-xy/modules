package scanner

import "time"

type Message struct {
    Id        uint64     `json:"id"`
    Content  []byte      `json:"content"`
    Bytes    int         `json:"bytes"`
    Timestamp time.Time  `json:"timestamp"`
}

func (m *Message) IsEmpty() bool {
    if m.Bytes == 0 {
        return true
    }

    if len(m.Content) == 0/* && len(m.Fields) == 0*/ {
        return true
    }

    return false
}
