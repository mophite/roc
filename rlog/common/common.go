package common

type Level uint8

const (
	NONE Level = 0x00 + iota
	DEBUG
	STACK
	INFO
	WARN
	ERR
	FATAL
)

type Config = string

type Detail struct {
	Name      string `json:"name"`
	Line      string `json:"line"`
	Prefix    string `json:"prefix"`
	Trace     uint64 `json:"trace"`
	Content   string `json:"content"`
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
}
