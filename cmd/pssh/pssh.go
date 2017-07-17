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

	jobs := make(chan psshutils.HostInfo, len(hosts))
	results := make(chan Result, len(hosts))

	var res = []Result{}

	// create workers
	for i := 1; i <= args.Parallel; i++ {
		go pssh(i, args.Command, args.Timeout, jobs, results)
	}

	// submit jobs
	for _, h := range hosts {
		jobs <- h
		log.Printf("submitting job host = %s", h.Addr)
	}

	// read results
	for i := 0; i < len(hosts); i++ {
		r := <- results
		log.Printf("reading results for host = %s", r.Addr)
		res = append(res, r)
	}
}

func pssh(worker int, command []string, timeout int, jobs <-chan psshutils.HostInfo, results chan<- Result) {
	for host := range jobs {
		stdout := make([]string, 1)
		stderr := make([]string, 1)

		log.Printf("pssh(%d) got job %s\n", worker, host.Addr)
		stdout[0] = "this is stdout for " + host.Addr
		stderr[0] = "no stderr"
		time.Sleep(1000000000)
		results <- Result{
			Stdout: stdout,
			Stderr: stderr,
			Code:   0,
			Addr:   host.Addr,
		}
	}
}
