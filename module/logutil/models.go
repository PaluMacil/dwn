package logutil

import (
	"fmt"
	"time"

	"github.com/PaluMacil/dwn/module/core"
)

type LogLevel int8

func (lvl LogLevel) String() string {
	return fmt.Sprintf("[%s]", lvl)
}

// none - no log message emitted
// error - Other runtime errors or unexpected conditions. Expect these to be immediately visible on a status console.
// warn - Use of deprecated APIs, poor use of API, 'almost' errors, other runtime situations that are undesirable or unexpected, but not necessarily "wrong". Expect these to be immediately visible on a status console.
// info - Interesting runtime events (startup/shutdown). Expect these to be immediately visible on a console, so be conservative and keep to a minimum.
// debug - detailed information on the flow through the system. Expect these to be written to logs only.
// trace - more detailed information. Expect these to be written to logs only.
const (
	LevelNone LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

const (
	ConfigPrefix = "LOGUTIL_CONFIG"
)

type Config struct {
	ConsoleLevel  LogLevel         `json:"consoleLevel"`
	UseColorCodes bool             `json:"useColorCodes"`
	QueueLevel    LogLevel         `json:"queueLevel"`
	ModifiedBy    core.DisplayName `json:"modifiedBy"`
	Modified      time.Time        `json:"modified"`
}

func (c Config) Key() []byte {
	return c.Prefix()
}

func (c Config) Prefix() []byte {
	return []byte(ConfigPrefix)
}

type Entry struct {
	Level   int
	Message string
}
