package flags

import (
	"flag"
	"time"
)

const (
	defaultClients  = 10
	defaultExecTime = 10 * time.Second

	defaultCmdTimeout      = 1 * time.Second
	defaultMaxThinkingTime = 100

	defaultKeySpaceSize = 10000
	defaultPayloadSize  = 1 << 8
)

type Flags struct {
	NumClients int
	ExecTime   time.Duration
	CmdTimeout time.Duration

	MaxThinkingTimeMs int
	KeySpaceSize      int
	PayloadSize       int
}

func ParseFlagsFromArgs() Flags {
	f := Flags{}

	flag.IntVar(&f.NumClients, "clients", defaultClients, "number of concurrent clients (int)")
	flag.DurationVar(&f.ExecTime, "exec-time", defaultExecTime, "total execution time (duration)")
	flag.DurationVar(&f.CmdTimeout, "cmd-timeout", defaultCmdTimeout, "command timeout (duration)")
	flag.IntVar(&f.MaxThinkingTimeMs, "max-thinking-time", defaultMaxThinkingTime, "maximum thinking time to wait between requests, in milliseconds (int)")
	flag.IntVar(&f.KeySpaceSize, "key-space", defaultKeySpaceSize, "number of different keys (int)")
	flag.IntVar(&f.PayloadSize, "payload-size", defaultPayloadSize, "payload size of values, in Bytes (int)")

	flag.Parse()
	return f
}
