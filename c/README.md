struct
======
- TODO(cjpatton) Function for computing BF parameters.

Implementations of (keyed) data structures and (keyed) hash functions:

  * `tiny.{h,c}` implements an efficient transformation of a cryptographic hash
    function (specifially SHA512) or pseudorandom function (specifically
    HMAC-SHA512) to new a function _H_ with a tiny range. It's designed to
    ensure that for each input _x_, the output _H(x)_ is uniform over the
    specified range. This makes it suitable for use in data structures specified
    here.

  * `bloom.{h,c}` implements a [Bloom
    filter](https://en.wikipedia.org/wiki/Bloom_filter), a data structure that
    concisely reprsents a set. Specifically, this implements the efficient
    _double-hashing scheme_ of [Kirsch and
    Mitzenmacher](https://www.eecs.harvard.edu/~michaelm/postscripts/rsa2008.pdf).
    The hash is instantiated with `tinyhash` or `tinyprf`.

 * `dict.{h,c}` implements a [Bloomier filter](https://arxiv.org/abs/0807.0928),
    a data structure that concisely represents a function from arbitrary strings
    to strings of some fixed length. Construction time is _expected_ O(_n_), where
    _n_ is the number of input/output pairs, and look up time is O(1). The hash
    functions are implemented using `hash`/`tinyhash` or `prf`/`tinyprf`.

 * `bits.{h,c}` implements some low-level bit-manipulation functions needed for
   everything else.

See `README.md` in the root project directory for installation instructions.
