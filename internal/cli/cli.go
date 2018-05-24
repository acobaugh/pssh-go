package cli

import (
	"os"
)

type CommonArgs struct {
	HostFiles []string `arg:"env:HOST_FILES,-f,separate"`
	Hosts     []string `arg:"env:HOSTS,-H,separate"`
	User      string   `arg:"-l"`
	Parallel  int      `arg:"-p"`
	Timeout   int      `arg:"-t"`
	ListHosts bool     `arg:"-L,help:List the hosts that were selected"`
	Verbose   bool     `arg:"-v"`
	Select    string   `arg:"-s,help:Shell-style glob to filter hosts"`
}

func (CommonArgs) Version() string {
	return os.Args[0] + " pssh-go 0.1"
}
