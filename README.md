# valpass - a small golang module to verify passwords 

## Background 

A decade ago I designed an encryption algorithm
just for fun and to learn more about cryptography.
During development I wrote a little helper tool
which I could use to verify some quality metrics
og my algorithm:
[analyze.c](https://github.com/TLINDEN/twenty4/blob/master/analyze/analyze.c).

This module is a re-implementation of this code
with go as a reusable module.

## Features

- standalone module without external dependencies
- uses 3 different metrics to measure password quality
- you can configure which metric to use
- you can also configure the quality thresholds
- there's support for dictionary lookup, but you need to provide the dictionary
- different metrics for ASCII and UTF-8 character space
- it's reasonably fast
- the code is small enough to just copy it into your code

## Quality metrics

![1000006662](https://github.com/user-attachments/assets/6cf19c6f-7c7a-4a2c-9a58-95b3ac1c49e7)

A good password is easy to remember and hard
to guess. Don't be fooled by those "use special characters"
evangelists: diceware passwords as outlined in the
well known xkcd comic are by far the best ones.

However, if it's your job zo implement a register
user interface, then sooner or later you'll need
to validate the password the user just entered.

This module can be used for this job.

By default it checks 3 metrics:

### Entropy

Entropy in this case measures the cryptographic
strength of the password. I non-technical words:
it checks how scrambled the password looks or how
many different bits it uses.

By default we only look for printable US-ASCII characters.

### Character diffusion

Of course just measuring entropy is insufficient. For
instance a password `12345678` consists of 8 different 
characters and might pass the entropy check. However, as
can be easily seen, the characters are sorted and 
therefore this password would be s terrible one.

Thus, character diffusion measures how characters are
distributed.

Keep in mind that these two metrics would flag
the `Tr0ub4dor&3` password of the comic as pretty good,
while in reality it's not! You might remedy 
this problem with a longer mandatory password 
length. But zhe harsh reality is, that people still 
use such passwords.

### Compression

We go one step further and also measure how much
the password can be compressed. For instance, let's 
look at this run length encoding example:

The string `aaabggthhhh` can be rle encoded to
`2ab2gt4h`. The ideal password cannot be compressed
or not much.

Of course ee do not use RLE. We measure compression 
using the [Flate algorithm](
https://en.m.wikipedia.org/wiki/Deflate).

### Optional: dictionary check

You can supply a dictionary of words of your
liking and check if the password under test
matches one if the words. Submatches can also 
be done.

### Custom

You can also enable or disable certain metrics and
you can tune the quality thresholds as needed.

### Future/ ToDo

- checksum test using supplied checksum list, e.g. of leaked passwords
- fuzzy testing against dictionary to catch variations
- chi square test (see  http://www.fourmilab.ch/random/)
- Arithmetic mean value test
- Monte Carlo value test
- Serial correlation
- maybe some dieharder tests


## Usage

Since the module is not yet complete and undocumented,
please look at [the example](https://github.com/TLINDEN/valpass/blob/main/example/test.go)
how to use it.

## License 

This module is licensed under the BSD license.

## Prior art


[go-password](https://github.com/wagslane/go-password-validator) provides similar
functionality and it's stable and battle tested. 
However ir only measures the character entropy.

