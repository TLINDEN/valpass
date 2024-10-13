package valpass

import (
	"bytes"
	"compress/flate"
	"fmt"
	"math"
	"strings"
)

/*
 * Contains the raw  dictionary data and some flags.  Must be provided
 * by the user
 */
type Dictionary struct {
	Words    []string // the actual dictionary
	Submatch bool     // if true 'foo' would match 'foobar'
}

/*
 * Options define how to operate the validation
 */
type Options struct {
	Compress         int         // minimum compression rate in percent
	CharDistribution float64     // minimum char distribution in percent
	Entropy          float64     // minimum entropy value in bits/char
	Dictionary       *Dictionary // if set, lookup given dictionary, the caller provides it
	UTF8             bool        // if true work on unicode utf-8 space, not just bytes
	Mean             float64     // if >0, calculate the arithmetic mean
}

/*
 * Default validation config, a compromise of comfort and security, as always.
 */
const (
	MIN_ENTROPY  float64 = 3.0
	MIN_COMPRESS int     = 10
	MIN_DICT     bool    = false
	MIN_DIST     float64 = 10.0
	MAX_UTF8     int     = 2164864 // max characters encodable with utf8
	MAX_CHARS    int     = 95      // maximum printable US ASCII chars
	MIN_DICT_LEN int     = 5000

	//  we start  our ascii  arrays  at char(32),  so to  have max  95
	// elements in the slice, we subtract 32 from each ascii code
	MIN_ASCII byte = 32

	//  arithmetic  mean limits:  we work on  chr(32) til  chr(126) in
	// ascii. The mean value, however, is not 63 as one would suppose,
	// but  80, because most used  printable ascii chars exist  in the
	// upper area  of the space. So,  we take 80 as  the middle ground
	// and go beyond 5 up or down
	MIDDLE_MEAN float64 = 80
	LIMIT_MEAN  float64 = 5
)

/*
Stores the results of all validations.
*/
type Result struct {
	Ok               bool    // overall result
	DictionaryMatch  bool    // true if the password matched a dictionary entry
	Compress         int     // actual compression rate in percent
	CharDistribution float64 // actual character distribution in percent
	Entropy          float64 // actual entropy value in bits/chars
	Mean             float64 // actual arithmetic mean, close to 127.5 is best
}

/*
 * Generic validation function. You should only call this function and
 * tune it  using the Options  struct. However, options  are optional,
 * there are sensible defaults builtin
 */
func Validate(passphrase string, opts ...Options) (Result, error) {
	result := Result{Ok: true}

	// defaults, see above
	options := Options{
		MIN_COMPRESS,
		MIN_DIST,
		MIN_ENTROPY,
		nil,
		false, // dict: default off
		0,     // mean: default off
	}

	if len(opts) == 1 {
		options = opts[0]
	}

	// execute the actual validation checks

	if options.Entropy > 0 {
		var entropy float64
		var err error

		switch options.UTF8 {
		case true:
			entropy, err = GetEntropyUTF8(passphrase)
			if err != nil {
				return result, err
			}
		default:
			entropy, err = GetEntropyAscii(passphrase)
			if err != nil {
				return result, err
			}
		}

		if entropy <= options.Entropy {
			result.Ok = false
		}

		result.Entropy = entropy
	}

	if options.Compress > 0 {
		compression, err := GetCompression([]byte(passphrase))
		if err != nil {
			return result, err
		}

		if compression >= options.Compress {
			result.Ok = false
		}

		result.Compress = compression
	}

	if options.CharDistribution > 0 {
		var dist float64

		switch options.UTF8 {
		case true:
			dist = GetDistributionUTF8(passphrase)
		default:
			dist = GetDistributionAscii(passphrase)
		}
		if dist <= options.CharDistribution {
			result.Ok = false
		}

		result.CharDistribution = dist
	}

	if options.Dictionary != nil {
		match, err := GetDictMatch(passphrase, options.Dictionary)
		if err != nil {
			return result, err
		}

		if match {
			result.Ok = false
			result.DictionaryMatch = true
		}
	}

	if options.Mean > 0 {
		mean := GetArithmeticMean(passphrase)

		if mean > (MIDDLE_MEAN+options.Mean) || mean < (MIDDLE_MEAN-options.Mean) {
			result.Ok = false
		}

		result.Mean = mean
	}

	return result, nil
}

/*
 * we  compress with  Flate level  9 (max)  and see  if the  result is
 * smaller than the password, in which case it could be compressed and
 * contains repeating characters;  OR it is larger  than the password,
 * in which case it could NOT be compressed, which is what we want.
 */
func GetCompression(passphrase []byte) (int, error) {
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
 * Return the  entropy as bits/rune, where  rune is a unicode  char in
 * utf8 space.
 */
func GetEntropyUTF8(passphrase string) (float64, error) {
	var entropy float64
	length := len(passphrase)

	wherechar := make([]int, MAX_UTF8)
	hist := make([]int, length)
	var histlen int

	for i := 0; i < MAX_UTF8; i++ {
		wherechar[i] = -1
	}

	for _, char := range passphrase {
		if wherechar[char] == -1 {
			wherechar[char] = histlen
			histlen++
		}

		hist[wherechar[char]]++
	}

	for i := 0; i < histlen; i++ {
		diff := float64(hist[i]) / float64(length)
		entropy -= diff * math.Log2(diff)
	}

	return entropy, nil
}

/*
Return the entropy as bits/char, where  char is a printable char in
US-ASCII space. Returns error if a char is non-printable.
*/
func GetEntropyAscii(passphrase string) (float64, error) {
	var entropy float64
	length := len(passphrase)

	wherechar := make([]int, MAX_CHARS)
	hist := make([]int, length)
	var histlen int

	for i := 0; i < MAX_CHARS; i++ {
		wherechar[i] = -1
	}

	for _, char := range []byte(passphrase) {
		if char < MIN_ASCII || char > 126 {
			return 0, fmt.Errorf("non-printable ASCII character encountered: %c", char)
		}
		if wherechar[char-MIN_ASCII] == -1 {
			wherechar[char-MIN_ASCII] = histlen
			histlen++
		}

		hist[wherechar[char-MIN_ASCII]]++
	}

	for i := 0; i < histlen; i++ {
		diff := float64(hist[i]) / float64(length)
		entropy -= diff * math.Log2(diff)
	}

	return entropy, nil
}

/*
 * Return character distribution in utf8 space
 */
func GetDistributionUTF8(passphrase string) float64 {
	hash := make([]int, MAX_UTF8)
	var chars float64

	for _, char := range passphrase {
		hash[char]++
	}

	for i := 0; i < MAX_UTF8; i++ {
		if hash[i] > 0 {
			chars++
		}
	}
	return chars / (float64(MAX_UTF8) / 100)
}

/*
 * Return character distribution in US-ASCII space
 */
func GetDistributionAscii(passphrase string) float64 {
	hash := make([]int, MAX_CHARS)
	var chars float64

	for _, char := range []byte(passphrase) {
		hash[char-MIN_ASCII]++
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
func GetDictMatch(passphrase string, dict *Dictionary) (bool, error) {
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

/*
* Return  the arithmetic  mean value:

	This is simply the result of summing the all the bytes (bits if the

-b  option  is specified)  in  the  file  and  dividing by  the  file
length. If the  data are close to random, this  should be about 127.5
(0.5 for -b option output). If  the mean departs from this value, the
values are consistently high or low.

	Working on US-ASCII space
*/
func GetArithmeticMean(passphrase string) float64 {
	sum := 0.0
	count := 0.0

	for _, char := range []byte(passphrase) {
		sum += float64(char)
		count++
	}

	return sum / count
}
