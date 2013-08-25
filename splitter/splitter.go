package splitter

import (
	"math/rand"

	"github.com/stuntgoat/snl/percent_sample"
)


type Splitter struct {
	Buckets [][]string

	// maps an index of `Buckets` that to a percentage of samples
	// for that bucket.
	PercentageMap map[int]int

	well []string
	count int
}

// distributeLines will distribute the values from the well into the Buckets
// using a probability distribution created from the
// objects PercentageMap values.
// This occurs when the program receives a SIGINT while scanning a file
// or if the file contains a number of lines that are not evenly
// divided by 100.
func (splitter *Splitter) DistributeLines() {
	var choice float64
	var val float64
	var probValue float64
	var line string

	// sort well
	percent_sample.Shuffle235(splitter.well, splitter.count)


	var totalValues = 0
	for _, percentVal := range splitter.PercentageMap {
		totalValues += percentVal
	}

	// for every value in the current well,
	// randomly select a bucket to place the value in the well
	for i := 0; i < splitter.count; i++ {
		choice = rand.Float64()
		val = 0.0

		for bucketIndex, intPercentValue := range splitter.PercentageMap {
			probValue = float64(intPercentValue) / totalValues

			val = val + probValue
			if choice < val {
				line = splitter.well[i]
				splitter.Buckets[bucketIndex] = append(splitter.Buckets[bucketIndex], line)
				break
			}
		}
	}
}

// AddLine places the line into the total sample
func (splitter *Splitter) AddLine(line string) {
	splitter.well = append(splitter.well, line)
	splitter.count++
}