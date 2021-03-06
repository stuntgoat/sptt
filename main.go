package main

import (
	"bufio"
	"fmt"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
	"github.com/stuntgoat/sptt/splitter"
)

const (
	TESTING_SET = iota
	TRAINING_SET
)

var FILENAME string
var PERCENTAGE_TRAIN int
var PERCENTAGE_TEST int
var command = os.Args[0]
var invocation = fmt.Sprintf("%s -train PERCENT FILE\n", command)
var invocationStdin = fmt.Sprintf("%s -train PERCENT -\n", command)

var SAMPLE = splitter.Splitter{}

var logger *log.Logger

// flag.Usage help message override
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage:\n%s%s", invocation, invocationStdin)
}

func init() {
	pTrainHelp := "Percentage of randomly selected lines for training dataset."
	ptrain := flag.Int("train", 0, pTrainHelp)

	logger = log.New(os.Stderr, "[sptt] ", log.LstdFlags|log.Lshortfile)
	rand.Seed(time.Now().UTC().UnixNano())

	flag.Usage = Usage
	flag.Parse()

	PERCENTAGE_TRAIN = *ptrain
	if PERCENTAGE_TRAIN < 0 || PERCENTAGE_TRAIN > 100 {
		logger.Println("PERCENT must be between 0 and 100")
		Usage()
		os.Exit(1)
	}
}





// parseFile validates a string and returns an *os.File
func parseFile(s string) (file *os.File) {
	if s == "" {
		logger.Print("[Error] missing filename argument")
		Usage()
		os.Exit(1)
	}

	file, err := os.Open(s)
	if err != nil {
		logger.Printf("[Error] error opening %s: %s", s, err)
		Usage()
		os.Exit(1)
	}

	return file
}

func writeLines() {
	trainingFile := FILENAME + ".train"
	testingFile := FILENAME + ".test"
	writeFiles(trainingFile, SAMPLE.Buckets[TRAINING_SET])
	writeFiles(testingFile, SAMPLE.Buckets[TESTING_SET])
}

// writeFiles writes the testing and training files to disk.
func writeFiles(name string, lines []string) {
	tf, err := os.Create(name)
	if err != nil {
		logger.Printf("[Error] unable to open %s: %s", name, err)
		os.Exit(1)
	}
	defer tf.Close()
	for _, line := range lines {
		_, err := tf.WriteString(line + "\n")
		if err != nil {
			logger.Print("[Error] unable to write line to test file")
			os.Exit(1)
		}
	}
}

// handleSignal handles a SIGINT (control-c) when the user
// might want to break from a stream while sampling a percentage.
func handleSignal() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	<- sigChannel

	SAMPLE.DistributeLines()
	writeLines()
	os.Exit(0)
}

func main() {
	var file *os.File
	var line string

	// calculate the number of samples that need to go into each bucket
	testBucket := make([]string, 0)
	trainBucket := make([]string, 0)
	PERCENTAGE_TEST = 100 - PERCENTAGE_TRAIN

	SAMPLE.Buckets = [][]string{testBucket, trainBucket}
	SAMPLE.PercentageMap = map[int]int{
		TESTING_SET: PERCENTAGE_TEST,
		TRAINING_SET: PERCENTAGE_TRAIN,
	}

	FILENAME = flag.Arg(0)
	if FILENAME == "-" {
		file = os.Stdin
		FILENAME = "STDIN"
	} else {
		file = parseFile(FILENAME)
		defer file.Close()
	}

	go handleSignal()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprint(scanner.Text())
		SAMPLE.AddLine(line)
	}
	SAMPLE.DistributeLines()
	writeLines()
}
