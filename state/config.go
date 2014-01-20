package state

import (
	"flag"
	"github.com/eelavai/octopus/exit"
	"math/rand"
	"strings"
	"stripe-ctf.com/log"
	"time"
)

type config struct {
	seed                                   int64
	nodeCount                              int
	root, containerRoot, sqlcluster, write string
	verbose, dryrun                        bool
	duration                               time.Duration
	containerIds, username                 string
	args                                   []string

	parsedContainerIds []string
	wg                 *exit.WaitGroup
}

var conf = &config{}

func AddFlags() {
	flag.Int64Var(&conf.seed, "seed", 0, "Seed for the PRNG. 0 means autoseed")
	flag.StringVar(&conf.root, "root", "/tmp/octopus",
		"The root directory, from Octopus's perspective")
	flag.StringVar(&conf.containerRoot, "container-root", "",
		"The root directory, from the perspective of each node.")
	flag.IntVar(&conf.nodeCount, "c", 3, "Number of servers to spawn")
	flag.BoolVar(&conf.dryrun, "n", false, "Just print out what commands would be run, don't actually run them")
	flag.BoolVar(&conf.verbose, "v", false, "Enable debug output")
	flag.StringVar(&conf.sqlcluster, "r", "", "Path to sqlcluster binary")
	flag.DurationVar(&conf.duration, "duration", 30*time.Second, "How long to run the octopus for")
	flag.StringVar(&conf.containerIds, "container-ids", "", "IDs of containers to run in")
	flag.StringVar(&conf.username, "username", "", "Username to run as")
	flag.StringVar(&conf.write, "w", "", "Where to write results")
}

func AfterParse() {
	conf.args = flag.Args()

	if conf.seed == 0 {
		conf.seed = randomSeed()
	}

	if conf.nodeCount < 2 {
		log.Fatal("You need at least two agents")
	}

	if conf.sqlcluster == "" {
		conf.sqlcluster = sqlclusterPath()
		if conf.sqlcluster == "" {
			log.Fatalf("No `sqlcluster' binary provided, and could not find one heuristically. Please provide a `-r <sqlcluster>' argument, run `%s/build.sh' to build the binary locally, or install it globally via `go install stripe-ctf.com'.", scriptDir())
		} else {
			log.Printf("Found sqlcluster binary at %s", conf.sqlcluster)
		}
	}

	if conf.containerRoot == "" {
		conf.containerRoot = conf.root
	}

	if conf.containerIds != "" {
		conf.parsedContainerIds = strings.Split(conf.containerIds, ",")
	}

	log.SetVerbose(conf.verbose)
	log.Printf("Using seed %v", conf.seed)

	// Always seed the default RNG. It controls comparatively little, so
	// don't worry about exposing it
	rand.Seed(randomSeed())
	conf.wg = exit.NewWaitGroup()

	recordConfig()
}

func Local() bool {
	return ContainerIds() == nil
}

func Seed() int64 {
	return conf.seed
}

func NodeCount() int {
	return conf.nodeCount
}

func Root() string {
	return conf.root
}

func ContainerRoot() string {
	return conf.containerRoot
}

func Sqlcluster() string {
	return conf.sqlcluster
}

func Write() string {
	return conf.write
}

func Verbose() bool {
	return conf.verbose
}

func Dryrun() bool {
	return conf.dryrun
}

func Duration() time.Duration {
	return conf.duration
}

func ContainerIds() []string {
	return conf.parsedContainerIds
}

func Username() string {
	return conf.username
}

func WaitGroup() *exit.WaitGroup {
	return conf.wg
}

func Args() []string {
	return conf.args
}
