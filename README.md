[![Go Report Card](https://goreportcard.com/badge/github.com/tlinden/valpass)](https://goreportcard.com/report/github.com/tlinden/valpass) 
[![Actions](https://github.com/tlinden/valpass/actions/workflows/ci.yaml/badge.svg)](https://github.com/tlinden/valpass/actions)
[![Go Coverage](https://github.com/tlinden/valpass/wiki/coverage.svg)](https://raw.githack.com/wiki/tlinden/valpass/coverage.html)
![GitHub License](https://img.shields.io/github/license/tlinden/valpass)
[![GoDoc](https://godoc.org/github.com/tlinden/valpass?status.svg)](https://godoc.org/github.com/tlinden/valpass)

# valpass - a small golang module to verify passwords 

## Background 

A decade ago I designed an encryption algorithm
just for fun and to learn more about cryptography.
During development I wrote a little helper tool
which I could use to verify some quality metrics
of my algorithm:
[analyze.c](https://github.com/TLINDEN/twenty4/blob/master/analyze/analyze.c).

This module is a re-implementation of this code
with go as a reusable module.

## Features

- standalone module without external dependencies
- uses 5 different metrics to measure password quality
- you can configure which metric to use
- you can also configure the quality thresholds
- there's support for dictionary lookup, but you need to provide the dictionary yourself 
- different metrics for ASCII and UTF-8 character space
- it's reasonably fast
- the code is small enough to just copy it into your code

## Quality metrics

![1000006662](https://github.com/user-attachments/assets/6cf19c6f-7c7a-4a2c-9a58-95b3ac1c49e7)

A good password is easy to remember and hard
to guess. Don't be fooled by those "use special characters"
evangelists: diceware passwords as outlined in the
well known xkcd comic are by far the best ones.

However, if it's your job to implement a registration 
user interface, then sooner or later you'll need
to validate passwords.

This module can be used for this job.

By default it checks 3 metrics:

### Entropy

Entropy in this case measures the cryptographic
strength of the password. In non-technical words:
it checks how scrambled the password looks or how
many different bits it uses.

By default we only look for printable US-ASCII characters. But you can switch to UTF-8 as well.

### Character diffusion

Of course just measuring entropy is insufficient. For
instance a password `12345678` consists of 8 different 
characters and might pass the entropy check. However, as
can be easily seen, the characters are sorted and 
therefore this password would be a terrible one.

Thus, character diffusion measures how characters are
distributed.

Keep in mind that these two metrics would flag
the `Tr0ub4dor&3` password of the comic as pretty good,
while in reality it's not! You might remedy 
this problem with a longer mandatory password 
length. But the harsh reality is that people still 
use such passwords.

### Compression

We go one step further and also measure how much
the password can be compressed. For instance, let's 
look at this run length encoding example:

The string `aaabggthhhh` can be rle encoded to
`2ab2gt4h`. The result is shorter than the original, it is compressed.
The ideal password cannot be compressed
or not much.

Of course we do not use RLE. We measure compression 
using the [Flate algorithm](
https://en.m.wikipedia.org/wiki/Deflate).

### Optional: arithmetic mean value

This is simply the result of summing the all the printable ascii chars
divided by password length. The ideal value would be ~80, because most
normal  letters hang  out in  the upper  area between  32 (space)  and
126(tilde). We  consider a password ok,  if its mean lies  around this
area give or  take 5.  If the  mean departs more from  this value, the
characters are consistently  high or low (e.g. more  numbers and upper
case  letters or  only  lower case  letters). The  latter,  5, can  be
tweaked. The larger the number, tha laxer the result.

Please be  warned, that this  metric will in  most cases give  you bad
results on otherwise good passwords,  such as diceware passwords. Only
use it if you know what you're doing.

### Optional: dictionary check

You can supply a dictionary of words of your
liking and check if the password under test
matches one of the words. Submatches can also 
be done.

### Custom measurements

You can also enable or disable certain metrics and
you can tune the quality thresholds as needed.

### Future/ ToDo

- checksum test using supplied checksum list, e.g. of leaked passwords
-  fuzzy  testing  against   dictionary  to  catch  variations,  using
  Levenshtein or something similar.


## Usage

Usage is pretty simple:

```go
import "github.com/tlinden/valpass"

[..]
   res, err := valpass.Validate("password"); if err != nil {
     log.Fatal(err)
   }
   
   if !res.Ok {
     log.Fatal("Password is unsecure!")
   }
[..]
```

You may also tune which tests you want to execute and with wich
parameters. To do this, just supply a second argument, which must be a
`valpas.Options` struct:

```go
type Options struct {
        Compress         int         // minimum compression rate in percent
        CharDistribution float64     // minimum char distribution in percent
        Entropy          float64     // minimum entropy value in bits/char
        Dictionary       *Dictionary // if set, lookup given dictionary, the caller provides it
        UTF8             bool        // if true work on unicode utf-8 space, not just bytes
}
```

To turn off a test, just set the tunable to zero.

Please take a look at [the
example](https://github.com/TLINDEN/valpass/blob/main/example/test.go)
or at [the unit tests](https://github.com/TLINDEN/valpass/blob/main/lib_test.go).

## Performance

Benchmark results of version 0.0.1:

```default
% go test -bench=. -count 5
goos: linux
goarch: amd64
pkg: github.com/tlinden/valpass
cpu: Intel(R) Core(TM) i7-10610U CPU @ 1.80GHz
BenchmarkValidateEntropy-8         98703             12402 ns/op
BenchmarkValidateEntropy-8         92745             12258 ns/op
BenchmarkValidateEntropy-8         94020             12495 ns/op
BenchmarkValidateEntropy-8         96747             12349 ns/op
BenchmarkValidateEntropy-8         94790             12368 ns/op
BenchmarkValidateCharDist-8        95610             12184 ns/op
BenchmarkValidateCharDist-8        96631             12305 ns/op
BenchmarkValidateCharDist-8        97537             12215 ns/op
BenchmarkValidateCharDist-8        97544             13703 ns/op
BenchmarkValidateCharDist-8        95139             15392 ns/op
BenchmarkValidateCompress-8         2140            636274 ns/op
BenchmarkValidateCompress-8         5883            204162 ns/op
BenchmarkValidateCompress-8         5341            229536 ns/op
BenchmarkValidateCompress-8         4590            221610 ns/op
BenchmarkValidateCompress-8         5889            186709 ns/op
BenchmarkValidateDict-8               81          13730450 ns/op
BenchmarkValidateDict-8               78          16081013 ns/op
BenchmarkValidateDict-8               74          17545981 ns/op
BenchmarkValidateDict-8               92          12830625 ns/op
BenchmarkValidateDict-8               94          12564205 ns/op
BenchmarkValidateAll-8              5084            200770 ns/op
BenchmarkValidateAll-8              6054            193329 ns/op
BenchmarkValidateAll-8              5998            186064 ns/op
BenchmarkValidateAll-8              5996            191017 ns/op
BenchmarkValidateAll-8              6268            173846 ns/op
BenchmarkValidateAllwDict-8          374           3054042 ns/op
BenchmarkValidateAllwDict-8          390           3109049 ns/op
BenchmarkValidateAllwDict-8          404           3022698 ns/op
BenchmarkValidateAllwDict-8          393           3075163 ns/op
BenchmarkValidateAllwDict-8          381           3112361 ns/op
PASS
ok      github.com/tlinden/valpass      54.017s
```

## License 

This module is licensed under the BSD license.

## Prior art


[go-password](https://github.com/wagslane/go-password-validator) provides similar
functionality and it's stable and battle tested. 
However ir only measures the character entropy.

