package pssh

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var DEFAULT_PORT = 22
var DEFAULT_USER = os.Getenv("USER")
var DEFAULT_PARALLEL = MaxParallelism()

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

type HostInfo struct {
	Addr string
	Port int
	User string
}

// Read in a series of host files and return a slice of HostInfo. Slice will be empty if any errors are encountered
func parseHostFileArgs(paths []string) []HostInfo {
	var hosts []HostInfo
	for _, f := range paths {
		h, err := readHostFile(f)
		if err != nil {
			continue
		}
		hosts = append(hosts, h...)
	}
	return hosts
}

// Read a host file and return a slice of HostInfo. error will contain any errors reading the file
func readHostFile(path string) ([]HostInfo, error) {
	var hosts []HostInfo

	file, err := os.Open(path)
	if err != nil {
		return hosts, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host, err := parseHostString(scanner.Text())
		if err != nil {
			log.Printf("readHostFile(): %s", err)
			continue
		}
		hosts = append(hosts, host)
	}
	return hosts, err
}

// Take a string of the form [user@]host[:port] and return a HostInfo struct
func parseHostString(str string) (HostInfo, error) {
	var err error

	re := regexp.MustCompile("((?P<user>.+)@)?(?P<addr>[a-zA-Z0-9.-]+)(:(?P<port>[0-9]+))?")
	names := re.SubexpNames()
	matches := re.FindAllStringSubmatch(str, -1)

	// if we didn't match, return an error
	if matches == nil {
		err = fmt.Errorf("Host string did not match regex: '%s'", str)
		return HostInfo{}, err
	}
	md := map[string]string{}

	// set md[name] = value of the named capture group
	for i, match := range matches[0] {
		md[names[i]] = match
	}

	// if we got an empty string or something that couldn't be turned into an int,
	// set port to the default port
	port, err := strconv.Atoi(md["port"])
	if err != nil {
		port = DEFAULT_PORT
	}

	// if we didn't get a user, set it to the default
	user := md["user"]
	if user == "" {
		user = DEFAULT_USER
	}

	return HostInfo{
		Addr: md["addr"],
		User: user,
		Port: port,
	}, nil
}

func parseHostStringArgs(hoststrings []string) []HostInfo {
	var hosts []HostInfo

	for _, v := range hoststrings {
		for _, s := range strings.Split(v, " ") {
			host, err := parseHostString(s)
			if err == nil {
				hosts = append(hosts, host)
			}
		}
	}
	return hosts
}

func GetHostsFromArgs(hostfiles []string, hoststrings []string) []HostInfo {
	return append(parseHostFileArgs(hostfiles), parseHostStringArgs(hoststrings)...)
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
