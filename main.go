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
var VALIDATION int
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

	validationHelp := "Number of validation files to create."
	validation := flag.Int("validation", 0, validationHelp)

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
	VALIDATION = *validation
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

func writeFiles() {
	trainingFile := FILENAME + ".train"
	testingFile := FILENAME + ".test"
	writeFile(trainingFile, SAMPLE.Buckets[TRAINING_SET])
	writeFile(testingFile, SAMPLE.Buckets[TESTING_SET])

	if VALIDATION == 0 {
		return
	}

	// create a new validation sample object from the training data bucket
	// use the number of buckets in VALIDATION
	vSample := splitter.Splitter{
		Buckets: make([][]string, 0),
		PercentageMap: make(map[int]int),
	}

	for i := 0; i < VALIDATION; i++ {
		vSample.Buckets = append(vSample.Buckets,  make([]string, 0))
	}

	vBucketPercent := int((100.0 / float64(VALIDATION)))
	for i := 0; i < VALIDATION; i++ {
		vSample.PercentageMap[i] = vBucketPercent
	}

	// call AddLine for each line in the training data
	for _, line := range SAMPLE.Buckets[TRAINING_SET] {
		vSample.AddLine(line)
	}

	vSample.DistributeLines()

	var nameWithSuffix string
	// call Write file for each bucket in the validation sample Bucket
	for i, bucket := range vSample.Buckets {
		nameWithSuffix = fmt.Sprintf("%s.V.%d", trainingFile, i + 1)
		writeFile(nameWithSuffix, bucket)
	}
}

// writeFiles writes the testing and training files to disk.
func writeFile(name string, lines []string) {
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
	writeFiles()
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

	// replace writeLines with writeFiles

	writeFiles()
}
