// Package valpass can be used to validate password quality using different metrics.
package valpass

import (
	"bytes"
	"compress/flate"
	"fmt"
	"math"
	"strings"
)

// Dictionary is a container struct to store and submit a dictionary of words.
type Dictionary struct {
	Words    []string // Contains the actual dictionary.
	Submatch bool     // Set to true to enable submatches, e.g. 'foo' would match 'foobar', default is false.
	Fuzzy    bool     // Set to true to enable more lax dictionary checks, default is false.
}

// Options struct can be used  to configure the validator, turn on/off
// certain validator functions and tune  the thresholds when to flag a
// password as valid.
//
// Set option to zero or false to disable the feature.
type Options struct {
	Compress         int         // minimum compression rate in percent, default 10%
	CharDistribution float64     // minimum character distribution in percent, default 10%
	Entropy          float64     // minimum entropy value in bits/char, default 3 bits/s
	Dictionary       *Dictionary // lookup given dictionary, the caller has to provide it
}

const (
	MIN_COMPRESS int     = 10
	MIN_DIST     float64 = 10.0
	MIN_ENTROPY  float64 = 3.0
	MIN_DICT_LEN int     = 5000
	MAX_CHARS    int     = 95 // maximum printable US ASCII chars

	//  we start  our ascii  arrays  at char(32),  so to  have max  95
	// elements in the slice, we subtract 32 from each ascii code
	ascii_base byte = 32
)

// Result stores the results of all validations.
type Result struct {
	Ok               bool    // overall result
	DictionaryMatch  bool    // true if the password matched a dictionary entry
	Compress         int     // actual compression rate in percent
	CharDistribution float64 // actual character distribution in percent
	Entropy          float64 // actual entropy value in bits/chars
}

// Validate  validates a given password.  You can  tune its  behavior
// using the Options struct. However,  options are optional, there are
// sensible defaults builtins.
//
// The returned Result struct returns the password quality.
func Validate(passphrase string, opts ...Options) (Result, error) {
	result := Result{Ok: true}

	// defaults, see above
	options := Options{
		Compress:         MIN_COMPRESS,
		CharDistribution: MIN_DIST,
		Entropy:          MIN_ENTROPY,
		Dictionary:       nil,
	}

	if len(opts) == 1 {
		options = opts[0]
	}

	// execute the actual validation checks

	if options.Entropy > 0 {
		var entropy float64
		var err error

		entropy, err = getEntropy(passphrase)
		if err != nil {
			return result, err
		}

		if entropy <= options.Entropy {
			result.Ok = false
		}

		result.Entropy = entropy
	}

	if options.Compress > 0 {
		compression, err := getCompression([]byte(passphrase))
		if err != nil {
			return result, err
		}

		if compression >= options.Compress {
			result.Ok = false
		}

		result.Compress = compression
	}

	if options.CharDistribution > 0 {
		var dist = getDistribution(passphrase)

		if dist <= options.CharDistribution {
			result.Ok = false
		}

		result.CharDistribution = dist
	}

	if options.Dictionary != nil {
		match, err := getDictMatch(passphrase, options.Dictionary)
		if err != nil {
			return result, err
		}

		if match {
			result.Ok = false
			result.DictionaryMatch = true
		}
	}

	return result, nil
}

/*
 * we  compress with  Flate level  9 (max)  and see  if the  result is
 * smaller than the password, in which case it could be compressed and
 * contains repeating characters;  OR it is larger  than the password,
 * in which case it could NOT be compressed, which is what we want.
 */
func getCompression(passphrase []byte) (int, error) {
	var b bytes.Buffer
	flater, _ := flate.NewWriter(&b, 9)

	if _, err := flater.Write(passphrase); err != nil {
		return 0, fmt.Errorf("failed to write to flate writer: %w", err)
	}

	if err := flater.Flush(); err != nil {
		return 0, fmt.Errorf("failed to flush flate writer: %w", err)
	}

	if err := flater.Close(); err != nil {
		return 0, fmt.Errorf("failed to close flate writer: %w", err)
	}

	// use floats to avoid division by zero panic
	length := float32(len(passphrase))
	compressed := float32(len(b.Bytes()))

	if compressed >= length {
		return 0, nil
	}

	percent := 100 - (compressed / (length / 100))

	return int(percent), nil
}

/*
Return the entropy as bits/char, where  char is a printable char in
US-ASCII space. Returns error if a char is non-printable.
*/
func getEntropy(passphrase string) (float64, error) {
	var entropy float64
	length := len(passphrase)

	wherechar := make([]int, MAX_CHARS)
	hist := make([]int, length)
	var histlen int

	for i := 0; i < MAX_CHARS; i++ {
		wherechar[i] = -1
	}

	for _, char := range []byte(passphrase) {
		if char < ascii_base || char > 126 {
			return 0, fmt.Errorf("non-printable ASCII character encountered: %c", char)
		}
		if wherechar[char-ascii_base] == -1 {
			wherechar[char-ascii_base] = histlen
			histlen++
		}

		hist[wherechar[char-ascii_base]]++
	}

	for i := 0; i < histlen; i++ {
		diff := float64(hist[i]) / float64(length)
		entropy -= diff * math.Log2(diff)
	}

	return entropy, nil
}

/*
 * Return character distribution in US-ASCII space
 */
func getDistribution(passphrase string) float64 {
	hash := make([]int, MAX_CHARS)
	var chars float64

	for _, char := range []byte(passphrase) {
		hash[char-ascii_base]++
	}

	for i := 0; i < MAX_CHARS; i++ {
		if hash[i] > 0 {
			chars++
		}
	}
	return chars / (float64(MAX_CHARS) / 100)
}

/*
 * Return true if password can be  found in given dictionary. This has
 * to be supplied by the user, we do NOT ship with a dictionary!
 */
func getDictMatch(passphrase string, dict *Dictionary) (bool, error) {
	if len(dict.Words) < MIN_DICT_LEN {
		return false, fmt.Errorf("provided dictionary is too small")
	}

	lcpass := strings.ToLower(passphrase)

	if dict.Submatch {
		for _, word := range dict.Words {
			if strings.Contains(strings.ToLower(word), lcpass) {
				return true, nil
			}
		}
	} else {
		for _, word := range dict.Words {
			if lcpass == strings.ToLower(word) {
				return true, nil
			}
		}
	}

	return false, nil
}
