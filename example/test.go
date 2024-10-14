package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/tlinden/valpass"
)

const template string = `
Metric                   Random         Threshhold  Result
------------------------------------------------------------------
Compression rate         0%%             min %d%%     %d%%
Character distribution   100%%           min %0.2f%%  %0.2f%%
Character entropy        8.0 bits/char  min %0.2f    %0.2f bits/char
Character redundancy     0.0%%           max %0.2f%%  %0.2f%%
Dictionary match         false          false       %t
------------------------------------------------------------------
Validation response                                 %t
`

func main() {
	opts := valpass.Options{
		Compress:         valpass.MIN_COMPRESS,
		CharDistribution: valpass.MIN_DIST,
		Entropy:          valpass.MIN_ENTROPY,
		Dictionary:       &valpass.Dictionary{Words: ReadDict("t/american-english")},
	}

	res, err := valpass.Validate(os.Args[1], opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		template,
		opts.Compress,
		res.Compress,
		opts.CharDistribution,
		res.CharDistribution,
		opts.Entropy,
		res.Entropy,
		100-opts.CharDistribution,
		100-res.CharDistribution,
		res.DictionaryMatch,
		res.Ok,
	)

	if !res.Ok {
		os.Exit(1)
	}
}

func ReadDict(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
