package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	//"golang.org/x/crypto/ssh"
	"github.com/cobaugh/pssh-go/psshutils"
	"log"
	"time"
)

type Result struct {
	Stdout []string
	Stderr []string
	Code   int
	Addr   string
	Ready  bool
}

func main() {
	var args struct {
		psshutils.CommonArgs
		Command []string `arg:"positional"`
	}

	args.Parallel = psshutils.DEFAULT_PARALLEL

	// parse args
	p := arg.MustParse(&args)
	if args.Verbose {
		fmt.Println("verbose defined")
	}

	// get hosts
	hosts := psshutils.GetHostsFromArgs(args.HostFile, args.HostString)

	// check arguments
	if len(args.Command) < 1 {
		p.Fail("No command specified")
	}
	if len(hosts) < 1 {
		p.Fail("No hosts specified")
	}

	inchan := make(chan psshutils.HostInfo, args.Parallel)
	outchan := make(chan Result)

	counter := 0
	var results = []Result{}

	for i := 1; i <= args.Parallel; i++ {
		go pssh(i, args.Command, args.Timeout, inchan, outchan)
	}

	for true {
		select {
		case inchan <- hosts[counter]:
			counter++
			if counter+1 >= len(hosts) {
				log.Printf("I am done, counter = %d\n", counter)
				close(inchan)
			} else {
				log.Printf("counter = %d, host = %s\n", counter, hosts[counter].Addr)
			}
		case r := <-outchan:
			log.Printf("reading results for host = %s", r.Addr)
			results = append(results, r)
		}
	}

}

func pssh(worker int, command []string, timeout int, inchan chan psshutils.HostInfo, outchan chan Result) {
	host := <- inchan
	stdout := make([]string, 1)
	stderr := make([]string, 1)

	log.Printf("pssh(%d) %s\n", worker, host.Addr)
	stdout[0] = "this is stdout for " + host.Addr
	stderr[0] = "no stderr"
	time.Sleep(1000000000)
	outchan <- Result{
		Stdout: stdout,
		Stderr: stderr,
		Code:   0,
		Addr:   host.Addr,
	}
}
