package message

import "time"

type Message struct {
    Magic     int                     `json:"magic"`
    Id        uint64                  `json:"id"`
    Header    map[string]interface{}  `json:"header"`
    Body      body                    `json:"body"`
    Footer    footer                  `json:"footer"`
}

type body struct {
    Content  []byte  `json:"content"`
    Bytes    int     `json:"bytes"`
}

type footer struct {
    CheckSum  int64  `json:"checksum"`
    Timestamp int64  `json:"timestamp"`
}

func New() *Message {
    return &Message{
        Magic: 0x0100,
        Id: 0,
        Header: map[string]interface{}{},
        Body:body{
            Bytes: -1,
        },
        Footer:footer{
            Timestamp: time.Now().Unix(),
        },
    }
}

func (m *Message) IsEmpty() bool {
	if m.Bytes == 0 {
		return true
	}

	if len(m.Content) == 0 && len(m.Fields) == 0 {
		return true
	}

	return false
}

func (m *Message) SetHeader(key, value string) error {
    return nil
}

func (m *Message) SetBody(body []byte) error {
    m.Body.Content = body
    m.Body.Bytes = len(body)
    return nil
}

func (m *Message) GetBodyLength() int64 {
    return int64(m.Body.Bytes)
}

func (m *Message) SetFooter() error {
    m.Footer.CheckSum = 0
    return nil
}
