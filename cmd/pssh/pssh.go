package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	//"golang.org/x/crypto/ssh"
	"log"
	"pssh-go/psshutils"
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

	numHosts := len(hosts)
	counter := 0
	var results = []Result{}

	go pssh(args.Command, args.Timeout, inchan, outchan)
	for true {
		select {
		case <-inchan:
			if counter-1 == numHosts {
				close(inchan)
			} else {
				inchan <- hosts[counter]
				counter++
			}
		case r := <-outchan:
			results = append(results, r)
		}
	}

}

func pssh(command []string, timeout int, inchan chan psshutils.HostInfo, outchan chan Result) {
	host := <-inchan

	stdout := make([]string, 1)
	stderr := make([]string, 1)

	log.Printf("pssh() %s\n", host.Addr)
	stdout[0] = "this is stdout for " + host.Addr
	stderr[0] = "no stderr"
	time.Sleep(5000000000)
	outchan <- Result{
		Stdout: stdout,
		Stderr: stderr,
		Code:   0,
		Addr:   host.Addr,
	}
}
