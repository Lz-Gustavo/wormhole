package flags

import (
	"flag"
	"strings"
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
	KeySpaceSize      int64
	PayloadSize       int

	LatencyMsrFilename string
	StatusMsrFilename  string
	EtcdHosts          []string
}

func ParseFlagsFromArgs() Flags {
	f := Flags{}

	flag.IntVar(&f.NumClients, "clients", defaultClients, "number of concurrent clients (int)")
	flag.DurationVar(&f.ExecTime, "exec-time", defaultExecTime, "total execution time (duration)")
	flag.DurationVar(&f.CmdTimeout, "cmd-timeout", defaultCmdTimeout, "command timeout (duration)")
	flag.IntVar(&f.MaxThinkingTimeMs, "max-thinking-time", defaultMaxThinkingTime, "maximum thinking time to wait between requests, in milliseconds (int)")
	flag.Int64Var(&f.KeySpaceSize, "key-space", defaultKeySpaceSize, "number of different keys (int64)")
	flag.IntVar(&f.PayloadSize, "payload-size", defaultPayloadSize, "payload size of values, in Bytes (int: 256, 512, 1024, 4096)")
	flag.StringVar(&f.LatencyMsrFilename, "latency-filename", "", "filename to write latency measurement, empty is disabled (string)")
	flag.StringVar(&f.StatusMsrFilename, "status-filename", "", "filename to write response status measurement, empty is disabled (string)")
	etcdHosts := flag.String("etcd-hosts", "", "list of etcd hostnames to send request, separated by `,` (string)")
	flag.Parse()

	if etcdHosts != nil {
		f.EtcdHosts = strings.Split(*etcdHosts, ",")
	}
	return f
}
