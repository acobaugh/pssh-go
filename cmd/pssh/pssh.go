package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	//"golang.org/x/crypto/ssh"
	"github.com/cobaugh/pssh-go/psshutils"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Result struct {
	Stdout string
	Stderr string
	Code   int
	Addr   string
	Ready  bool
}

func main() {
	var args struct {
		psshutils.CommonArgs
		Command string `arg:"positional"`
	}

	args.Parallel = psshutils.DEFAULT_PARALLEL

	// parse args
	p := arg.MustParse(&args)
	if args.Verbose {
		fmt.Println("verbose defined")
	}

	// get hosts
	hosts := psshutils.GetHostsFromArgs(args.HostFiles, args.Hosts)

	// check arguments
	if args.Command == "" {
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
		r := <-results
		log.Printf("reading results for host = %s, stdout = %s, stderr = %s", r.Addr, r.Stdout, r.Stderr)
		res = append(res, r)
	}
}

func pssh(worker int, command string, timeout int, jobs <-chan psshutils.HostInfo, results chan<- Result) {
	for host := range jobs {
		log.Printf("pssh(%d) got job %s@%s:%d\n", worker, host.User, host.Addr, host.Port)

		stderr := ""

		cmd := exec.Command("ssh", "-l", host.User, "-p", strconv.Itoa(host.Port), host.Addr, command)
		printCommand(cmd)
		stdout, err := cmd.Output()
		if err != nil {
			stderr = err.Error()
		}

		results <- Result{
			Stdout: string(stdout[:]),
			Stderr: stderr,
			Addr:   host.Addr,
			Code:   0,
		}
	}
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}
