package event

import "github.com/rookie-xy/modules/agents/file/scanner"

type Event struct {
    Magic     int              `json:"magic"`
    Header    Header           `json:"header"`
    Body     *scanner.Message  `json:"body"`
    Footer    Footer           `json:"footer"`
}

type Header map[string]interface{}

type Footer struct {
    CheckSum  []byte  `json:"checksum"`
}
